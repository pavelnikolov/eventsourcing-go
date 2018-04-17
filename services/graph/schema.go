package graph

// Schema is the GraphQL schema of the service
var Schema = `
	schema {
		query: Query
	}
	
	# The query type, represents all of the entry points into our object graph
	type Query {
		# article queries for an article by the provided id.
		article(id: ID!): Article
		# articles queries for latest artciles by category and status. If category is not provided it returns latest articles from all categories. 
		articles(category: String, count: Int! = 10, status: ArticleStatus! = PUBLISHED): [Article]!
	}

	enum ArticleStatus {
		UNKNOWN
		DRAFT
		PUBLISHED
		RETRACTED
	}

	type Article {
		id: ID!
		title: String!
		body: String!
		category: String!
		author_id: ID!
		author_name: String!
		status: ArticleStatus!
	}
`
