package submission

import (
 "testing"

 "github.com/Ploos-AS/Ploform/backend/internal/formschema"
)

func TestContactSubmission(t *testing.T){
 schema:=formschema.Contact()
 good:=map[string][]string{"name":{" Ada "},"email":{"ada@example.com"},"message":{"Hei!"}}
 got,err:=Validate(schema,good)
 if err!=nil{t.Fatal(err)}
 if got["name"]!="Ada"||got["email"]!="ada@example.com"||got["message"]!="Hei!"{t.Fatalf("unexpected normalized values: %#v",got)}
 if _,ok:=got["subject"];ok{t.Fatal("optional empty subject should be omitted")}
}

func TestContactRejectsInvalidValues(t *testing.T){
 cases:=[]struct{name string;values map[string][]string}{
  {"missing name",map[string][]string{"email":{"a@example.com"},"message":{"Hi"}}},
  {"bad email",map[string][]string{"name":{"Ada"},"email":{"bad"},"message":{"Hi"}}},
  {"missing message",map[string][]string{"name":{"Ada"},"email":{"a@example.com"}}},
  {"unknown field",map[string][]string{"name":{"Ada"},"email":{"a@example.com"},"message":{"Hi"},"admin":{"true"}}},
  {"duplicate field",map[string][]string{"name":{"Ada","Bob"},"email":{"a@example.com"},"message":{"Hi"}}},
  {"oversize",map[string][]string{"name":{"Ada"},"email":{"a@example.com"},"message":{string(make([]byte,10001))}}},
 }
 for _,tc:=range cases{t.Run(tc.name,func(t *testing.T){if _,err:=Validate(formschema.Contact(),tc.values);err==nil{t.Fatal("expected validation error")}})}
}

func TestSelectAndCheckbox(t *testing.T){
 schema:=formschema.Schema{Fields:[]formschema.Field{
  {ID:"agree",Type:"checkbox",Label:"Agree",Required:true},
  {ID:"color",Type:"select",Label:"Color",Options:[]string{"red","blue"}},
 }}
 if _,err:=Validate(schema,map[string][]string{"agree":{"on"},"color":{"green"}});err==nil{t.Fatal("unknown option accepted")}
 if _,err:=Validate(schema,map[string][]string{"agree":{"on"},"color":{"red"}});err!=nil{t.Fatal(err)}
 if _,err:=Validate(schema,map[string][]string{"color":{"red"}});err==nil{t.Fatal("required checkbox missing")}
}
