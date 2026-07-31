---
menuTitle: Route ordering
title: Route ordering in delegation
weight: 10
description: Understand how Gloo Gateway orders routes when using VirtualService delegation, and how to control that order with weights.
---

When Gloo Gateway processes a delegated routing tree, it flattens the routes from all matched RouteTables into a single list before sending them to Envoy. How those routes are ordered in the final list depends on whether you use direct references or selector-based delegation, and whether RouteTables share the same weight. This page explains the ordering rules and how to control them.

## Inline routes and direct RouteTable references

When a VirtualService defines routes inline, or delegates to a single RouteTable via `delegateAction.ref`, routes are passed to Envoy in the exact order they are defined. No automatic reordering occurs.

Envoy evaluates routes top-to-bottom and uses the first match, so you are responsible for placing more specific routes before broader catch-all routes.

## Selector-based delegation

When a VirtualService delegates traffic to RouteTables via the `delegateAction.selector`, Gloo Gateway matches all RouteTables whose labels satisfy the selector and merges their routes. How the merged routes are ordered depends on the `weight` values that are assigned to each RouteTable.

### Distinct weights

When each matched RouteTable has a unique `weight` value, Gloo Gateway processes RouteTables in ascending weight order. Lower weights are assigned higher priority. Routes within each RouteTable appear in the final Envoy list exactly as you wrote them in the YAML, from top to bottom. When building the merged list, Gloo Gateway concatenates the RouteTables in weight order without sorting individual routes across RouteTable boundaries. This means that no route from a higher-weight RouteTable can appear before a route from a lower-weight RouteTable.

The `weight` field defaults to `0` when not set. 

In the following example, the routes from the `rt-specific` RouteTable have a higher priority than routes from the `rt-broad` RouteTable. 

{{< highlight yaml "hl_lines=8 25" >}}
apiVersion: gateway.solo.io/v1
kind: RouteTable
metadata:
  name: rt-specific
  labels:
    app: myapp
spec:
  weight: 10          # processed first (lowest weight)
  routes:
    - matchers:
        - regex: /api/homepage-beta-flag.*
      routeAction:
        single:
          upstream:
            name: beta-upstream
            namespace: gloo-system
---
apiVersion: gateway.solo.io/v1
kind: RouteTable
metadata:
  name: rt-broad
  labels:
    app: myapp
spec:
  weight: 20          # processed second
  routes:
    - matchers:
        - regex: /api/homepage.*
      routeAction:
        single:
          upstream:
            name: main-upstream
            namespace: gloo-system
{{< /highlight >}}

### Same weight

When two or more RouteTables that are matched by the same selector share the same `weight` (including the default of `0`), Gloo Gateway merges their routes into a single list and applies an automatic path-specificity sort before sending the result to Envoy.

The sort applies the following precedence rules, in order:

| Priority | Path type | Tiebreaker |
|----------|-----------|------------|
| 1 (highest) | exact | — |
| 2 | prefix | Longer prefix wins |
| 3 (lowest) | regex | Alphabetical comparison of pattern text |

For exact and prefix matches, the sort reliably places more specific routes before broader ones. For regex matches, meaningful specificity ordering is not possible. A regex pattern's text does not reflect how many or how few requests it matches, so the sort falls back to comparing the pattern text alphabetically, character by character using ASCII values. 

**Example ordering of regex patterns**

Let's assume you define the following regex patterns in your RouteTable: 
* `/(?i)api/homepage-beta-flag.*`
* `/(?i)api/homepage.*`

The two patterns are identical up to `homepage`. Then, the patterns are different with `-` as the next character in the first pattern and `.` in the second. Because `.` has a higher ASCII value than `-`, `/(?i)api/homepage.*` is placed before `/(?i)api/homepage-beta-flag.*` in the sorted list. Envoy evaluates `/(?i)api/homepage.*` before `/(?i)api/homepage-beta-flag.*`, and the more specific beta-flag pattern might never be reached. To ensure the more specific pattern is evaluated first, assign the RouteTable that contains it a lower weight.

{{% notice warning %}}
The automatic sort is triggered whenever two or more RouteTables share the same weight, including the default weight of `0`. Adding a new RouteTable with the same labels and no explicit weight causes Gloo Gateway to re-sort the merged routes from all same-weight RouteTables. Routes that were previously in user-defined order might move. To prevent resorting, assign each RouteTable a unique weight.
{{% /notice %}}
