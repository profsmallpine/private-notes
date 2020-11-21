package domain

// import (
// 	"gorm.io/gorm"
// )

// const (
// 	QueryJoinKey        = "join"
// 	QueryOffsetLimitKey = "offset-limit"
// 	QueryPreloadKey     = "preload"
// 	QuerySortKey        = "sort"
// 	QueryWhereKey       = "where"
// )

// // Queryable is the interface a type implements in order to be leveraged in a QueryService helper function.
// type Queryable interface {
// 	// NOTE: make sure the implementation is defined on a pointer receiver to have the
// 	// compiler force us to pass pointers into our QueryService methods.
// 	TableName() string
// }

// // QueryService is the interface describing queries that an integration with a database must support.
// type QueryService interface {
// 	// NOTE(dlk): Prefer QueryService methods over the gorm.DB API for reliable behavior.
// 	RawDB() *gorm.DB

// 	// NOTE(joe): The queryable arg should always be a pointer. We can't define this at the interface level
// 	// because *Queryable indicates a pointer to the definition of the Queryable interface.
// 	CountByQuery(queryable Queryable, query map[string]interface{}) (int64, error)
// 	DeleteByID(queryable Queryable, ID interface{}) error
// 	DeleteByQuery(queryable Queryable, query map[string]interface{}) error
// 	// NOTE(dlk): FetchByQuery's listPointer has to be a pointer to a slice, example: *[]*domain.User
// 	FetchByQuery(queryable Queryable, listPointer interface{}, query map[string]interface{}) error
// 	FindByID(queryable Queryable, ID interface{}) error
// 	FindByQuery(queryable Queryable, query map[string]interface{}) error
// 	FindOrInsert(queryable Queryable, findQueryable Queryable) error
// 	Insert(queryable Queryable) error
// 	LastByQuery(queryable Queryable, query map[string]interface{}) error
// 	UpdateAllByQuery(queryable Queryable, query, updates map[string]interface{}) error
// 	UpdateByID(queryable Queryable, ID interface{}, updates map[string]interface{}) error
// 	UpdateByStruct(queryablePointer Queryable) error

// 	// NOTE(joe): These are for managing PG advisory locks. These are used to ensure single execution of things
// 	// across servers (e.g., data migrations, cron jobs).
// 	// Here are some helpful links on advisory locks:
// 	//   https://www.postgresql.org/docs/9.2/view-pg-locks.html
// 	//   https://hashrocket.com/blog/posts/advisory-locks-in-postgres, https://github.com/heroku/pg_lock (Heroku Ruby implementation)
// 	ObtainLock(lockKey uint32) (bool, error)
// 	ReleaseLock(lockKey uint32) error
// }

// // QueryJoin structures JOIN statements between tables.
// //
// // BUG(dlk): do not use when building queries resulting in UPDATE or DELETE statements.
// type QueryJoin struct {
// 	Clause string
// 	Params []interface{}
// }

// // A QueryOffsetLimit structures OFFSET ... LIMIT ... clauses. Both fields are not required.
// type QueryOffsetLimit struct {
// 	Limit  int
// 	Offset int
// }

// // A QueryPreload structures requests to preload assocations between structs, leveraging GORM's preloading
// // functionality.
// type QueryPreload struct {
// 	Unscope             bool
// 	Preloads            []string
// 	ConditionalPreloads map[string][]interface{}
// }

// // A QuerySort structures ORDER BY clauses.
// type QuerySort struct {
// 	// Named UnsafeClause because for other Query components, the structs encourage/enforce
// 	// parameterization (and so does Gorm). But for an order by clause, this isn't the case so
// 	// it is on the developer to make sure SQL injection is not possible.
// 	// Here is a closed issue in Gorm's repo where the maintainer essentially pushes this
// 	// responsibility onto the consumer: https://github.com/jinzhu/gorm/pull/1000
// 	UnsafeClause string
// }

// // A QueryWhere structures WHERE clauses replacing raw map[string]interface{} query building when = operator is not
// // desired; e.g., IN, >, !=.
// type QueryWhere struct {
// 	Clauses []string
// 	Params  []interface{}
// 	Unscope bool
// }
