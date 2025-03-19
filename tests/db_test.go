package tests

import (
	"errors"
	"testing"
	"vectorDb/db"
	"vectorDb/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDbInitialization(t *testing.T) {

	t.Run("testing the initialization of the db with Client and Store ", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		// mockClient:=client.NewMockClient(controller)
		mockClient:=mock.NewMockClient(controller)
		mockStore:=mock.NewMockStore(controller)
		db:=db.NewVectorDbWithClientAndStore(mockClient,mockStore)
		assert.NotEqual(t,nil,db)
	})

	
	t.Run("testing the initialization of the db only with Client ", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		// mockClient:=client.NewMockClient(controller)
		mockClient:=mock.NewMockClient(controller)
		// mockStore:=mock.NewMockStore(controller)
		db:=db.NewVectorDbWithClient(mockClient)
		assert.NotEqual(t,nil,db)
	})

	t.Run("testing the initialization of the db only with Store ", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		// mockClient:=client.NewMockClient(controller)
		// mockClient:=mock.NewMockClient(controller)
		mockStore:=mock.NewMockStore(controller)
		db:=db.NewVectorDbWithStore(mockStore)
		assert.NotEqual(t,nil,db)
	})

}

func TestDbSave(t *testing.T) {

	t.Run("testing the save method of db when it returns an error ", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		// mockClient:=client.NewMockClient(controller)
		mockClient:=mock.NewMockClient(controller)
		mockStore:=mock.NewMockStore(controller)
		db:=db.NewVectorDbWithClientAndStore(mockClient,mockStore)
		assert.NotEqual(t,nil,db)

		testStore:="testStore"
		mockStore.EXPECT().Save(testStore).Return(errors.New("error during save"))
		err:=db.Save(testStore)
		assert.NotEqual(t,nil,err)
		require.Equal(t,err.Error(),"error during save")
	})

	t.Run("testing the save method of db when it does not return an error ", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		// mockClient:=client.NewMockClient(controller)
		mockClient:=mock.NewMockClient(controller)
		mockStore:=mock.NewMockStore(controller)
		db:=db.NewVectorDbWithClientAndStore(mockClient,mockStore)
		assert.NotEqual(t,nil,db)

		testStore:="testStore"
		mockStore.EXPECT().Save(testStore).Return(nil)
		err:=db.Save(testStore)
		assert.Equal(t,nil,err)
	})
}

func TestDbLoad(t *testing.T) {

	t.Run("testing the load method of db when it returns an error ", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		// mockClient:=client.NewMockClient(controller)
		mockClient:=mock.NewMockClient(controller)
		mockStore:=mock.NewMockStore(controller)
		db:=db.NewVectorDbWithClientAndStore(mockClient,mockStore)
		assert.NotEqual(t,nil,db)

		testStore:="testStore"
		mockStore.EXPECT().Load(testStore).Return(errors.New("error during load"))
		err:=db.Load(testStore)
		assert.NotEqual(t,nil,err)
		require.Equal(t,err.Error(),"error during load")
	})

	t.Run("testing the load method of db when it does not return an error ", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		// mockClient:=client.NewMockClient(controller)
		mockClient:=mock.NewMockClient(controller)
		mockStore:=mock.NewMockStore(controller)
		db:=db.NewVectorDbWithClientAndStore(mockClient,mockStore)
		assert.NotEqual(t,nil,db)

		testStore:="testStore"
		mockStore.EXPECT().Load(testStore).Return(nil)
		err:=db.Load(testStore)
		assert.Equal(t,nil,err)
	})
}


func TestDbInsert(t *testing.T) {

	t.Run("testing the insert method of db when client.embed fails ", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		// mockClient:=client.NewMockClient(controller)
		mockClient:=mock.NewMockClient(controller)
		mockStore:=mock.NewMockStore(controller)
		db:=db.NewVectorDbWithClientAndStore(mockClient,mockStore)
		assert.NotEqual(t,nil,db)

		testKey:="testKey"
		testStore:="testStore"
		// embedding:=generateRandomFloat64Array(8)
		mockClient.EXPECT().Embed(testKey).Return(nil,errors.New("error during insert"))
		err:=db.Insert(testStore,testKey)
		assert.NotEqual(t,nil,err)
		require.Equal(t,err.Error(),"error during insert")
	})

	t.Run("testing the insert method of db when client.embed succeeds but store.insert fails ", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		// mockClient:=client.NewMockClient(controller)
		mockClient:=mock.NewMockClient(controller)
		mockStore:=mock.NewMockStore(controller)
		db:=db.NewVectorDbWithClientAndStore(mockClient,mockStore)
		assert.NotEqual(t,nil,db)

		testKey:="testKey"
		testStore:="testStore"
		embedding:=generateRandomFloat32Array(8)
		mockClient.EXPECT().Embed(testKey).Return(embedding,nil)
		mockStore.EXPECT().Insert(testStore,embedding,testKey).Return(errors.New("inserting in store failed"))
		err:=db.Insert(testStore,testKey)
		assert.NotEqual(t,nil,err)
		require.Equal(t,err.Error(),"inserting in store failed")
	})

	
	t.Run("testing the insert method of db when client.embed succeeds and store.insert succeeds ", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		// mockClient:=client.NewMockClient(controller)
		mockClient:=mock.NewMockClient(controller)
		mockStore:=mock.NewMockStore(controller)
		db:=db.NewVectorDbWithClientAndStore(mockClient,mockStore)
		assert.NotEqual(t,nil,db)

		testKey:="testKey"
		testStore:="testStore"
		embedding:=generateRandomFloat32Array(8)
		mockClient.EXPECT().Embed(testKey).Return(embedding,nil)
		mockStore.EXPECT().Insert(testStore,(embedding),testKey).Return(nil)
		err:=db.Insert(testStore,testKey)
		assert.Equal(t,nil,err)
	})

}


// TODO : Implement testing for Lookup and Search