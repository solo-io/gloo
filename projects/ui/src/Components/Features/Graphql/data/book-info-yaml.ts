export const bookInfoYsml = `apiVersion: graphql.gloo.solo.io/v1alpha1
kind: GraphQLSchema
metadata:
 creationTimestamp: "2022-01-20T05:17:07Z"
  generation: 1
  name: bookinfo-graphql
  namespace: gloo-system
  resourceVersion: "28779"
  uid: fbd7f8c1-45f3-4e7f-a7a4-a839d6fb623f
spec:
  executableSchema:
    executor:
      local:
        enableIntrospection: true
        resolutions:
          Query|productsForHome:
            restResolver:
              request:
                headers:
                  :method: GET
                  :path: /api/v1/products
              upstreamRef:
                name: default-productpage-9080
                namespace: gloo-system
          author:
            restResolver:
              request:
                headers:
                  :method: GET
                  :path: /details/{$parent.id}
              response:
                resultRoot: author
              upstreamRef:
                name: default-details-9080
                namespace: gloo-system
          pages:
            restResolver:
              request:
                headers:
                  :method: GET
                  :path: /details/{$parent.id}
              response:
                resultRoot: pages
              upstreamRef:
                name: default-details-9080
                namespace: gloo-system
          ratings:
            restResolver:
              request:
                headers:
                  :method: GET
                  :path: /ratings/{$parent.id}
              response:
                resultRoot: ratings[*]
                setters:
                  numStars: '[*][1]'
                  reviewer: '[*][0]'
              upstreamRef:
                name: default-ratings-9080
                namespace: gloo-system
          reviews:
            restResolver:
              request:
                headers:
                  :method: GET
                  :path: /reviews/{$parent.id}
              response:
                resultRoot: reviews
              upstreamRef:
                name: default-reviews-9080
                namespace: gloo-system
          year:
            restResolver:
              request:
                headers:
                  :method: GET
                  :path: /details/{$parent.id}
              response:
                resultRoot: year
              upstreamRef:
                name: default-details-9080
                namespace: gloo-system
    schema_definition: |
      type Query {
        productsForHome: [Product] @resolve(name: "Query|productsForHome")
      }

      type Product {
        id: String
        title: String
        descriptionHtml: String
        author: String @resolve(name: "author")
        pages: Int @resolve(name: "pages")
        year: Int @resolve(name: "year")
        reviews : [Review] @resolve(name: "reviews")
        ratings : [Rating] @resolve(name: "ratings")
      }

      type Review {
        reviewer: String
        text: String
      }

      type Rating {
        reviewer : String
        numStars : Int
      }
`;
