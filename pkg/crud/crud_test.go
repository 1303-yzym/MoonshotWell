package crud_test

//import (
//	"context"
//	"testing"
//	"time"
//
//	"github.com/1303-yzym/MoonshotWell/pkg/contract"
//	"gorm.io/gorm"
//
//	"github.com/stretchr/testify/assert"
//
//)

//var (
//	Storage        contract.Storage
//	CRUDRepository CRUDRepo
//)
//
//func TestMain(m *testing.M) {
//	Storage = testmock.NewUnitTestStorage()
//	CRUDRepository = NewCRUDRepo(Storage)
//
//	m.Run()
//}
//
//type CRUDRepo interface {
//	repository.CRUD[sqls.CrudTestModel]
//}
//
//type CRUDRepoImpl struct {
//	contract.Storage
//	repository.Crud[sqls.CrudTestModel]
//}
//
//func NewCRUDRepo(storage contract.Storage) CRUDRepo {
//	return &CRUDRepoImpl{
//		Storage: storage,
//		Crud:    repository.New[sqls.CrudTestModel](storage),
//	}
//}
//
//var testCrudModel = &sqls.CrudTestModel{
//	Model:        sqls.Model{Id: 100},
//	Str:          "dada",
//	Uid:          7334164863024893952,
//	Boolean:      ty.True(),
//	TsTime:       ty.NewTime(23, 21, 32),
//	TsData:       ty.NewData(2021, 6, 1),
//	TsDatatime:   ty.NewUnixTime(time.Now()),
//	JsonStrArray: ty.StrArray{"da1", "da2", "c"},
//	JsonMap:      ty.BMap{"dad": "da", "da2": ty.StrArray{"da1", "da2", "c"}},
//	JsonStrMap:   ty.StrMap{"da": "dad", "da2": "dd"},
//	TasksBit:     ty.Bitmap[[16]byte]{},
//	TaskCounters: ty.SMap[int64, int64]{
//		1: 1,
//	},
//}
//
//func TestCrudAdd(t *testing.T) {
//	err := CRUDRepository.Add(context.Background(), testCrudModel)
//	if err != nil {
//		t.Error(err)
//	}
//
//	assert.NotZero(t, testCrudModel.Id)
//}
//
//func TestCrudDelete(t *testing.T) {
//	ctx := context.Background()
//
//	err := CRUDRepository.DeleteById(ctx, 100, false)
//	if err != nil {
//		t.Error(err)
//	}
//	// err = CRUDRepository.DeleteById(ctx, 100, true)
//	// if err != nil {
//	// 	t.Error(err)
//	// }
//}
//
//func TestCrudType(t *testing.T) {
//	ctx := context.Background()
//
//	model, err := CRUDRepository.QueryById(ctx, 100)
//	if err != nil {
//		t.Error(err)
//	}
//
//	t.Log(model.TasksBit.String())
//	model.TasksBit.Set(4, true)
//	t.Log(model.TasksBit.String())
//
//	err = CRUDRepository.Update(ctx, func(db *gorm.Statement) {
//		db.Where("id = ?", model.Id)
//	}, "tasks_bit", model.TasksBit)
//	if err != nil {
//		t.Error(err)
//	}
//}
//
//func TestCrudUpdate(t *testing.T) {
//	err := CRUDRepository.Updates(context.Background(), func(db *gorm.Statement) {
//		db.Where("id = ?", testCrudModel.Id)
//	}, sqls.CrudTestModel{
//		Boolean: ty.False(),
//	})
//	if err != nil {
//		t.Error(err)
//	}
//
//	err = CRUDRepository.Update(context.Background(), func(db *gorm.Statement) {
//		db.Where("id = ?", testCrudModel.Id)
//	}, "ts_data", ty.NewData(2021, 10, 1))
//	if err != nil {
//		t.Error(err)
//	}
//}
//
//func TestCrudQuery2(t *testing.T) {
//	ctx := context.Background()
//
//	model, err := CRUDRepository.QueryById(ctx, testCrudModel.Id)
//	if err != nil {
//		t.Error(err)
//	}
//
//	assert.Equal(t, testCrudModel.TaskCounters, model.TaskCounters)
//}
//
//func TestCrudQuery(t *testing.T) {
//	ctx := context.Background()
//
//	model, err := CRUDRepository.QueryById(ctx, testCrudModel.Id)
//	if err != nil {
//		t.Error(err)
//	}
//
//	assert.Equal(t, testCrudModel.JsonStrArray, model.JsonStrArray)
//
//	first, err := CRUDRepository.QueryFirst(ctx, func(db *gorm.Statement) {
//		db.Where("uid = ?", testCrudModel.Uid)
//	})
//	if err != nil {
//		return
//	}
//
//	assert.Equal(t, testCrudModel.JsonStrArray, first.JsonStrArray)
//
//	firstNoErr := CRUDRepository.QueryFirstNoErr(ctx, func(db *gorm.Statement) {
//		db.Where("uid = ?", testCrudModel.Uid)
//	})
//	assert.NotNil(t, firstNoErr)
//	assert.Equal(t, testCrudModel.JsonStrArray, firstNoErr.JsonStrArray)
//
//	list, count, err := CRUDRepository.QueryPagination(ctx, ty.Pagination{
//		Page:     1,
//		PageSize: 10,
//	}, func(db *gorm.Statement) {
//		db.Where("str = dada").Order("id desc")
//	})
//	if err != nil {
//		return
//	}
//
//	assert.Equal(t, len(list), 2)
//	assert.Equal(t, count, 2)
//	assert.Equal(t, list[1].Id, testCrudModel.Id)
//}
//
//func TestCrudTransaction(t *testing.T) {
//	err := Storage.DB().Debug().Transaction(func(tx *gorm.DB) error {
//		err := CRUDRepository.WithDB(tx).Add(context.Background(), testCrudModel)
//		if err != nil {
//			return err
//		}
//
//		err = CRUDRepository.WithDB(tx).Updates(context.Background(), func(db *gorm.Statement) {
//			db.Where("id = ?", testCrudModel.Id)
//		}, sqls.CrudTestModel{
//			Boolean: ty.False(),
//		})
//		if err != nil {
//			return err
//		}
//
//		return nil
//	})
//	if err != nil {
//		t.Error(err)
//
//		return
//	}
//}
