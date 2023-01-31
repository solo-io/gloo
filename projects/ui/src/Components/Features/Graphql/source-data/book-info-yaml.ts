export const bookInfoYaml = `apiVersion: graphql.gloo.solo.io/v1beta1
kind: GraphQLApi
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"graphql.gloo.solo.io/v1beta1","kind":"GraphQLApi","metadata":{"annotations":{},"name":"bookinfo-graphql","namespace":"gloo-system"},"spec":{"executableSchema":{"executor":{"local":{"enableIntrospection":true,"resolutions":{"Query|productsForHome":{"restResolver":{"request":{"headers":{":method":"GET",":path":"/api/v1/products"}},"upstreamRef":{"name":"default-productpage-9080","namespace":"gloo-system"}}},"author":{"restResolver":{"request":{"headers":{":method":"GET",":path":"/details/{$parent.id}"}},"response":{"resultRoot":"author"},"upstreamRef":{"name":"default-details-9080","namespace":"gloo-system"}}},"pages":{"restResolver":{"request":{"headers":{":method":"GET",":path":"/details/{$parent.id}"}},"response":{"resultRoot":"pages"},"upstreamRef":{"name":"default-details-9080","namespace":"gloo-system"}}},"ratings":{"restResolver":{"request":{"headers":{":method":"GET",":path":"/ratings/{$parent.id}"}},"response":{"resultRoot":"ratings[*]","setters":{"numStars":"[*][1]","reviewer":"[*][0]"}},"upstreamRef":{"name":"default-ratings-9080","namespace":"gloo-system"}}},"reviews":{"restResolver":{"request":{"headers":{":method":"GET",":path":"/reviews/{$parent.id}"}},"response":{"resultRoot":"reviews"},"upstreamRef":{"name":"default-reviews-9080","namespace":"gloo-system"}}},"year":{"restResolver":{"request":{"headers":{":method":"GET",":path":"/details/{$parent.id}"}},"response":{"resultRoot":"year"},"upstreamRef":{"name":"default-details-9080","namespace":"gloo-system"}}}}}},"schema_definition":"type Query {\n  productsForHome: [Product] @resolve(name: \"Query|productsForHome\")\n}\n\ntype Product {\n  id: String\n  title: String\n  descriptionHtml: String\n  author: String @resolve(name: \"author\")\n  pages: Int @resolve(name: \"pages\")\n  year: Int @resolve(name: \"year\")\n  reviews : [Review] @resolve(name: \"reviews\")\n  ratings : [Rating] @resolve(name: \"ratings\")\n}\n\ntype Review {\n  reviewer: String\n  text: String\n}\n\ntype Rating {\n  reviewer : String\n  numStars : Int\n}\n"}}}
  creationTimestamp: "2022-01-25T05:21:31Z"
  generation: 1
  managedFields:
  - apiVersion: graphql.gloo.solo.io/v1beta1
    fieldsType: FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          .: {}
          f:kubectl.kubernetes.io/last-applied-configuration: {}
      f:spec:
        .: {}
        f:executableSchema:
          .: {}
          f:executor:
            .: {}
            f:local:
              .: {}
              f:enableIntrospection: {}
              f:resolutions:
                .: {}
                f:Query|productsForHome:
                  .: {}
                  f:restResolver:
                    .: {}
                    f:request:
                      .: {}
                      f:headers:
                        .: {}
                        f::method: {}
                        f::path: {}
                    f:upstreamRef:
                      .: {}
                      f:name: {}
                      f:namespace: {}
                f:author:
                  .: {}
                  f:restResolver:
                    .: {}
                    f:request:
                      .: {}
                      f:headers:
                        .: {}
                        f::method: {}
                        f::path: {}
                    f:response:
                      .: {}
                      f:resultRoot: {}
                    f:upstreamRef:
                      .: {}
                      f:name: {}
                      f:namespace: {}
                f:pages:
                  .: {}
                  f:restResolver:
                    .: {}
                    f:request:
                      .: {}
                      f:headers:
                        .: {}
                        f::method: {}
                        f::path: {}
                    f:response:
                      .: {}
                      f:resultRoot: {}
                    f:upstreamRef:
                      .: {}
                      f:name: {}
                      f:namespace: {}
                f:ratings:
                  .: {}
                  f:restResolver:
                    .: {}
                    f:request:
                      .: {}
                      f:headers:
                        .: {}
                        f::method: {}
                        f::path: {}
                    f:response:
                      .: {}
                      f:resultRoot: {}
                      f:setters:
                        .: {}
                        f:numStars: {}
                        f:reviewer: {}
                    f:upstreamRef:
                      .: {}
                      f:name: {}
                      f:namespace: {}
                f:reviews:
                  .: {}
                  f:restResolver:
                    .: {}
                    f:request:
                      .: {}
                      f:headers:
                        .: {}
                        f::method: {}
                        f::path: {}
                    f:response:
                      .: {}
                      f:resultRoot: {}
                    f:upstreamRef:
                      .: {}
                      f:name: {}
                      f:namespace: {}
                f:year:
                  .: {}
                  f:restResolver:
                    .: {}
                    f:request:
                      .: {}
                      f:headers:
                        .: {}
                        f::method: {}
                        f::path: {}
                    f:response:
                      .: {}
                      f:resultRoot: {}
                    f:upstreamRef:
                      .: {}
                      f:name: {}
                      f:namespace: {}
          f:schema_definition: {}
    manager: kubectl-client-side-apply
    operation: Update
    time: "2022-01-25T05:21:31Z"
  name: bookinfo-graphql
  namespace: gloo-system
  resourceVersion: "2405"
  uid: ca63a276-ab64-46e1-8aef-fa18d0a68157
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
