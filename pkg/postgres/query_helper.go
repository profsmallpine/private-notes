package postgres

// import (
// 	"strings"

// 	"github.com/profsmallpine/private-notes/domain"
// 	"gorm.io/gorm"
// )

// // buildQuery passes the query map through the supported clauses thereby constructing a query for execution.
// func buildQuery(querySoFar *gorm.DB, originalQuery map[string]interface{}) (newQuery *gorm.DB) {
// 	// NOTE(joe): Copy originalQuery map, otherwise when we delete keys from it, we are mutating
// 	// the original map passed into the QueryService. This means a map variable re-used
// 	// could be different when used for a 2nd query. By copying we avoid that possibility.
// 	queryCopy := map[string]interface{}{}
// 	for k, v := range originalQuery {
// 		queryCopy[k] = v
// 	}
// 	querySoFar = addJoin(querySoFar, queryCopy)
// 	querySoFar = querySoFar.Where(queryCopy)
// 	querySoFar = addOffsetLimit(querySoFar, queryCopy)
// 	querySoFar = addPreloads(querySoFar, queryCopy)
// 	querySoFar = addRawWhere(querySoFar, queryCopy)
// 	querySoFar = addSort(querySoFar, queryCopy)
// 	return querySoFar
// }

// // addJoin converts a domain.QueryJoin into a GORM query and adds it to the query under construction.
// func addJoin(querySoFar *gorm.DB, query map[string]interface{}) (newQuery *gorm.DB) {
// 	if _, ok := query[domain.QueryJoinKey]; ok {
// 		joins := query[domain.QueryJoinKey].([]domain.QueryJoin)
// 		for _, join := range joins {
// 			querySoFar = querySoFar.Joins(join.Clause, join.Params...)
// 		}
// 		delete(query, domain.QueryJoinKey)
// 	}

// 	return querySoFar
// }

// // addOffsetLimit converts a domain.QueryOffsetLimit into a GORM query and adds it to the query under
// // construction. If either field is it's zero-value, it is not included in the query.
// func addOffsetLimit(querySoFar *gorm.DB, query map[string]interface{}) (newQuery *gorm.DB) {
// 	ol, ok := query[domain.QueryOffsetLimitKey].(domain.QueryOffsetLimit)
// 	if ok {
// 		if ol.Limit > 0 {
// 			querySoFar = querySoFar.Limit(ol.Limit)
// 		}
// 		if ol.Offset > 0 {
// 			querySoFar = querySoFar.Offset(ol.Offset)
// 		}
// 		delete(query, domain.QueryOffsetLimitKey)
// 	}

// 	return querySoFar
// }

// // addPreloads converts a domain.QueryPreloads into a GORM query and adds it to the query under construction.
// //
// // TODO(dlk): consider updates to preloads with GORM 2.0
// func addPreloads(querySoFar *gorm.DB, query map[string]interface{}) (newQuery *gorm.DB) {
// 	if _, ok := query[domain.QueryPreloadKey]; ok {
// 		preload := query[domain.QueryPreloadKey].(domain.QueryPreload)
// 		for _, load := range preload.Preloads {
// 			// NOTE(joe): If we end up needing to have scoped and unscoped preloads, then we need
// 			// 			 to add a map to separate that out.
// 			conditions, conditionsExist := preload.ConditionalPreloads[load]

// 			if preload.Unscope {
// 				// NOTE(joe): GORM by default filters out soft-deleted records (deleted_at NOT NULL)
// 				querySoFar = querySoFar.Preload(load, func(db *gorm.DB) *gorm.DB {
// 					if conditionsExist {
// 						return db.Unscoped().Where(conditions[0], conditions[1:]...)
// 					} else {
// 						return db.Unscoped()
// 					}
// 				})
// 				continue
// 			}

// 			if conditionsExist {
// 				querySoFar = querySoFar.Preload(load, conditions...)
// 			} else {
// 				querySoFar = querySoFar.Preload(load)
// 			}
// 		}
// 		delete(query, domain.QueryPreloadKey)
// 	}

// 	return querySoFar
// }

// // addRawWhere converts a domain.QueryWhere into a GORM query and adds it to the query under construction.
// func addRawWhere(querySoFar *gorm.DB, query map[string]interface{}) (newQuery *gorm.DB) {
// 	if _, ok := query[domain.QueryWhereKey]; ok {
// 		where := query[domain.QueryWhereKey].(domain.QueryWhere)

// 		wrappedClauses := []string{}
// 		for _, c := range where.Clauses {
// 			wrappedClauses = append(wrappedClauses, "("+c+")")
// 		}
// 		clause := strings.Join(wrappedClauses, " AND ")

// 		if where.Unscope {
// 			querySoFar = querySoFar.Unscoped().Where(clause, where.Params...)
// 		} else {
// 			querySoFar = querySoFar.Where(clause, where.Params...)
// 		}
// 		delete(query, domain.QueryWhereKey)
// 	}

// 	return querySoFar
// }

// // addSort converts a domain.QuerySort into a GORM query and adds it to the query under construction.
// func addSort(querySoFar *gorm.DB, query map[string]interface{}) (newQuery *gorm.DB) {
// 	if _, ok := query[domain.QuerySortKey]; ok {
// 		sort := query[domain.QuerySortKey].(domain.QuerySort)
// 		querySoFar = querySoFar.Order(sort.UnsafeClause)
// 		delete(query, domain.QuerySortKey)
// 	}

// 	return querySoFar
// }
