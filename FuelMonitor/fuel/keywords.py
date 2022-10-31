from fuel.methods import makeQuery


def MonitorByKeyword(keyword: str):
    query = """query monitorByKeyword($keyword: String!){
  products(search: $keyword, pageSize: 30, sort: {updated_at: DESC}) {
    items {
      name
      sku
      id
      price {
          regularPrice {
              amount {
                  value
              }
          }
      }
      small_image {
        url
      }
      created_at
      updated_at
      stock_status
      __typename
      ... on ConfigurableProduct {
        variants {
            product {
                id
                created_at
                updated_at
                id 
                name
                sku
                stock_status
                is_raffle_item
            }
        }
      }
    }
  }
}
"""
    response = makeQuery(query, {'keyword': keyword})
    if response.success:
        return response.data
    else:
        print("Unexpected error: " + response.data)
