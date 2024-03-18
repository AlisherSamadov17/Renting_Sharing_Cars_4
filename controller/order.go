package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rent-car/models"
	"github.com/google/uuid"
)




func (c Controller) Order(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.CreateOrder(w,r)
	case http.MethodGet:
		values := r.URL.Query()
		_, ok := values["id"]
		if !ok {
			c.GetAllOrder(w, r)
		} else {
			c.GetByIDOrder(w,r)
		} 
	   case http.MethodPut:
		values := r.URL.Query()
		_, ok := values["id"]
		if ok {
			c.UpdateOrder(w, r)
		}
	// case http.MethodDelete:
		// values := r.URL.Query()
		// _, ok := values["id"]
		// if ok {
		// 	c.DeleteCustomer(w, r)
		// }
	default:
		handleResponse(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (c Controller) CreateOrder(w http.ResponseWriter,r *http.Request)  {
	order:=models.CreateOrder{}

	if err := json.NewDecoder(r.Body).Decode(&order);err != nil {
		errStr := fmt.Sprintf("error while decoding request body,err:%v\n",err)
		fmt.Println(errStr)
		handleResponse(w,http.StatusBadRequest,errStr)
		return
	}
	
	id,err :=c.Store.Order().Create(order)
	if err != nil{
		fmt.Println("error while creating order,err:",err)
		handleResponse(w,http.StatusInternalServerError,err)
		return
	}
	handleResponse(w,http.StatusOK,id)
}

func (c Controller) UpdateOrder(w http.ResponseWriter,r *http.Request)  {
	order:=models.UpdateOrder{}
	if err := json.NewDecoder(r.Body).Decode(&order);err != nil{
		errStr := fmt.Sprintf("error while decoding request body,err:%v\n",err)
		fmt.Println(errStr)
		handleResponse(w,http.StatusBadRequest,errStr)
		return
	}
	order.Id = r.URL.Query().Get("id")
	err := uuid.Validate(order.Id)
	if err != nil {
		fmt.Println("error while validating,err",err)
		handleResponse(w,http.StatusBadRequest,err.Error())
	return
	}
	
id,err :=c.Store.Order().Update(order)
if err != nil {
	fmt.Println("error while updating customer,err",err)
	handleResponse(w,http.StatusInternalServerError,err)
	return
}
handleResponse(w,http.StatusOK,id)
}

func (c Controller) GetAllOrder(w http.ResponseWriter,r *http.Request)  {
	var (
		values  = r.URL.Query()
		search  string
		request = models.GetAllOrdersRequest{}
	)
	if _, ok := values["search"]; ok {
		search = values["search"][0]
	}

	request.Search = search

	page, err := ParsePageQueryParam(r)
	if err != nil {
		fmt.Println("error while parsing page, err: ", err)
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	limit, err := ParseLimitQueryParam(r)
	if err != nil {
		fmt.Println("error while parsing limit, err: ", err)
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("page: ", page)
	fmt.Println("limit: ", limit)

	request.Page = page
	request.Limit = limit
	orders, err := c.Store.Order().GetAll(request)
	if err != nil {
		fmt.Println("error while getting orders, err: ", err)
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, orders)
}

func (c Controller) GetByIDOrder(w http.ResponseWriter,r *http.Request)  {
	values:= r.URL.Query()
	id :=values["id"][0]

	order,err := c.Store.Order().GetByID(id)
	if err != nil {
		fmt.Println("error while getting order by id:",err)
		handleResponse(w,http.StatusInternalServerError,err)
		return
	}
	handleResponse(w,http.StatusOK,order)
}


func (c Controller) DeleteOrder(w http.ResponseWriter,r *http.Request)  {
	id := r.URL.Query().Get("id")
	fmt.Println("id",id)

	err := uuid.Validate(id)
	if err != nil {
		fmt.Println("error while validating id,err:",err.Error())
	handleResponse(w,http.StatusBadRequest,err.Error())
	return
}
err = c.Store.Customer().Delete(id)
if err != nil {
	fmt.Println("error while deleting customer, err:",err)
	handleResponse(w,http.StatusInternalServerError,err)
return
}
handleResponse(w,http.StatusOK,id)
}