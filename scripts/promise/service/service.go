package promises

import (
	"github.com/satori/go.uuid"
)

const saltSize = 16

type Promise struct{
	dal *Dal
	name string
}

func CreatePromiseService() *Promise{
	promise := new(Promise)
	promise.name = "Promise"
	promise.dal = CreateDal(promise.name)
	
	return promise
}

func (f *Promise) AddPromise(userId uuid.UUID, title string, description string){
	f.dal.AddPromise(userId, title, description)
}

func (f *Promise) EditPromise(promiseId int, title string, description string, promisedTo uuid.UUID, status Status, privacy Privacy){
	f.dal.EditPromise(promiseId, title, description, promisedTo, status, privacy)
}

func (f *Promise) DeletePromise(promiseId int){
	f.dal.DeletePromise(promiseId)
}

func (f *Promise) GetPromises(userId uuid.UUID) []LEPromise{
	return f.dal.GetPromises(userId)
}

func (f *Promise) GetPromise(promiseId int) LEPromise{
	return f.dal.GetPromise(promiseId)
}

func (f *Promise) AskForPromise(userId uuid.UUID, promiseId int){
	
}

func (f *Promise) MakePromise(userId uuid.UUID, promiseId int){

}

func (f *Promise) Close(){
	f.dal.CloseConnection()
}

