# Gloo Helm Charts with Read-only UI

This directory contains templates for the values, Chart, and requirements files
which produce a Helm Chart for Gloo that includes a read-only version of the UI.

The read-only property is enforced by the apiserver (checks presence of a valid license).

For cleanliness, we use helm values to restrict the apiserver's role's permissions to
just those needed for read-only operations.

There are no helm templates unique to this chart. All templates are either used as
dependencies (Gloo's chart) or copied from the GlooE templates dir during build (apiserver templates).
