# Product Catalog Service

The Product Catalog Service governs the items catalog, search functions, categories, and item prices.

## What it does
* **Catalog Management**: Handles additions, deletions, and edits of product listings.
* **Product Search**: Searches catalog items by keyword matches.
* **Internal Catalog Validation**: Exposes internal lookup API used by Cart and Order services to verify item prices and existences.

## Data Flow
* **Public Reads**: 
  `Client -> API Gateway (bypasses auth checks) -> Product Catalog Service (GET /api/v1/products/)`
* **Internal Inquiries**: 
  `ShoppingCart Service -> Product Service (GET /internal/products/{productId})`

## Context Passing
* **User Identity Context**: Readers bypass authentication. Writes (POST/PUT/DELETE) require valid Admin/User context verified via gateway credentials mapping.

## Covered Aspects
* Keyword string matching queries in product titles and descriptions.
* In-memory dataset representing inventory catalogs.
* Price verification layers checking accurate double-precision pricing bounds.

## Future Aspects
* **Advanced Search (Elasticsearch)**: Integrate Elasticsearch or Algolia to provide typo tolerance, auto-completion, and facets.
* **Reviews and Ratings**: Add user review ratings database linking product IDs to user feedback.
* **Recommendation System**: Introduce catalog recommendation algorithms showing related products to users.
