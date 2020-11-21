package postgres

// import (
// 	"fmt"
// 	"reflect"
// 	"time"

// 	"github.com/profsmallpine/private-notes/domain"
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/clause"
// )

// // A QueryService holds a gorm.DB connection enabling querying the database.
// type QueryService struct {
// 	db *gorm.DB
// }

// // NewQueryService constructs a new *QueryService from the provided gorm.DB connection.
// func NewQueryService(DB *gorm.DB) *QueryService {
// 	return &QueryService{db: DB}
// }

// // RawDB exposes the gorm.DB instance for working with that API directly.
// func (qs *QueryService) RawDB() *gorm.DB {
// 	return qs.db
// }

// // FindByID gets the first record for the domain.Queryable matching the ID. ID must be a string, uint, or int.
// func (qs *QueryService) FindByID(queryable domain.Queryable, ID interface{}) error {
// 	qs.validateType(ID, []string{"string", "uint", "int"})
// 	return qs.db.First(queryable, ID).Error
// }

// // FetchByQuery gets all records for the domain.queryable matching the quest, storing all results in list. List must
// // be a pointer to a data type capable of storing many records; e.g., &[]*domain.User.
// func (qs *QueryService) FetchByQuery(queryable domain.Queryable, list interface{}, query map[string]interface{}) error {
// 	qs.validateType(list, []string{"ptr"})
// 	builtQuery := buildQuery(qs.db.Table(queryable.TableName()), query)
// 	return builtQuery.Find(list).Error
// }

// // FindByQuery gets the first record for the domain.Queryable matching the query.
// func (qs *QueryService) FindByQuery(queryable domain.Queryable, query map[string]interface{}) error {
// 	builtQuery := buildQuery(qs.db, query)
// 	return builtQuery.First(queryable).Error
// }

// // LastByQuery gets the last record for the domain.Queryable matching the query.
// func (qs *QueryService) LastByQuery(queryable domain.Queryable, query map[string]interface{}) error {
// 	builtQuery := buildQuery(qs.db, query)
// 	return builtQuery.Last(queryable).Error
// }

// // Insert creates a record for the domain.Queryable.
// func (qs *QueryService) Insert(queryable domain.Queryable) error {
// 	// NOTE(dlk): do not update Assocations
// 	return qs.db.Omit(clause.Associations).Create(queryable).Error
// }

// // FirstOrInsert gets the first record matching the fields of the second passed in domain.Queryable, filling the first
// // domain.Queryable with the resulting data, or, finding none, creates a record for the second passed in
// // domain.Queryable, filling the first domain.Queryable with the newly created data.
// func (qs *QueryService) FindOrInsert(queryable domain.Queryable, findQueryable domain.Queryable) error {
// 	return qs.db.FirstOrCreate(queryable, findQueryable).Error
// }

// // UpdateByID updates columns matching keys in the map for the first record of the domain.Queryable matching the ID.
// // ID must be a string, uint, or int.
// func (qs *QueryService) UpdateByID(queryable domain.Queryable, ID interface{}, updates map[string]interface{}) error {
// 	if _, ok := updates["updated_at"]; !ok {
// 		updates["updated_at"] = time.Now()
// 	}
// 	qs.validateType(ID, []string{"string", "uint", "int"})
// 	return qs.db.Table(queryable.TableName()).Where("id = ?", ID).Updates(updates).Error
// }

// // UpdateByStruct updates all non-zero value fields of the domain.Queryable.
// func (qs *QueryService) UpdateByStruct(queryablePointer domain.Queryable) error {
// 	// Make sure if an association is preloaded, we don't update that record too
// 	return qs.db.Omit(clause.Associations).Save(queryablePointer).Error
// }

// // UpdateAllByQuery updates columns matching keys in the map for all records of the domain.Queryable matching the query.
// func (qs *QueryService) UpdateAllByQuery(queryable domain.Queryable, query, updates map[string]interface{}) error {
// 	if _, ok := updates["updated_at"]; !ok {
// 		updates["updated_at"] = time.Now()
// 	}
// 	builtQuery := buildQuery(qs.db.Table(queryable.TableName()), query)
// 	return builtQuery.Updates(updates).Error
// }

// // DeleteByID performs a soft delete (i.e., populates "deleted_at") on the first record of the domain.Queryable matching
// // the ID.
// //
// // DeleteByID - and soft deleting in general - does not respect ON DELETE constraints.
// func (qs *QueryService) DeleteByID(queryable domain.Queryable, ID interface{}) error {
// 	qs.validateType(ID, []string{"string", "uint", "int"})
// 	return qs.db.Where("id = ?", ID).Delete(queryable).Error
// }

// // DeleteByQuery performs a soft delete (i.e., populates "deleted_at") on all records of the domain.Queryable matching
// // the query.
// func (qs *QueryService) DeleteByQuery(queryable domain.Queryable, query map[string]interface{}) error {
// 	builtQuery := buildQuery(qs.db.Table(queryable.TableName()), query)
// 	return builtQuery.Delete(queryable).Error
// }

// // CountByQuery performs a count of records for the domain.Queryable matching the query.
// func (qs *QueryService) CountByQuery(queryable domain.Queryable, query map[string]interface{}) (int64, error) {
// 	var count int64
// 	builtQuery := buildQuery(qs.db.Table(queryable.TableName()), query)

// 	//NOTE(joe): GORM doesn't add default scope to the Count method
// 	builtQuery = builtQuery.Where(queryable.TableName() + ".deleted_at IS NULL")
// 	err := builtQuery.Count(&count).Error
// 	return count, err
// }

// type advisoryLock struct {
// 	Obtained bool
// }

// func (qs *QueryService) ObtainLock(lockKey uint32) (bool, error) {
// 	var l advisoryLock
// 	r := qs.db.Raw("SELECT pg_try_advisory_lock(?) AS obtained;", lockKey)
// 	r.Scan(&l)
// 	return l.Obtained, r.Error
// }

// func (qs *QueryService) ReleaseLock(lockKey uint32) error {
// 	return qs.db.Exec("SELECT pg_advisory_unlock(?);", lockKey).Error
// }

// func (qs *QueryService) validateType(obj interface{}, validTypes []string) {
// 	kind := reflect.TypeOf(obj).Kind().String()

// 	for _, t := range validTypes {
// 		if kind == t {
// 			return
// 		}
// 	}

// 	issue := fmt.Sprintf("[FATAL ERROR] Argument must be one of: %v. Got %s instead.", validTypes, kind)
// 	fmt.Println(issue)
// 	panic(issue)
// }
