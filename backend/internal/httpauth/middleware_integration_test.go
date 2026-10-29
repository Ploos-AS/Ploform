package httpauth

import (
 "context"
 "net/http"
 "net/http/httptest"
 "testing"

 "github.com/Ploos-AS/Ploform/backend/internal/authz"
 "github.com/Ploos-AS/Ploform/backend/internal/session"
)

type testSession struct {identity session.Identity; err error}
func (s testSession) Lookup(context.Context,string)(session.Identity,error){return s.identity,s.err}
type testMembership struct {roles map[string]authz.Role}
func (m testMembership) Resolve(_ context.Context,_ string,org string)(authz.Role,error){
 role,ok:=m.roles[org];if !ok{return "",authz.ErrNoMembership};return role,nil
}

func TestMiddlewareSessionAndOrganizationIsolation(t *testing.T){
 const orgA="00000000-0000-0000-0000-000000000001"
 const orgB="00000000-0000-0000-0000-000000000002"
 const token="valid-session"
 tests:=[]struct{name string; identity session.Identity; sessionErr error; roles map[string]authz.Role; requestedOrg string; cookie bool; want int}{
  {"valid member",session.Identity{UserID:"user-a"},nil,map[string]authz.Role{orgA:authz.Owner},orgA,true,http.StatusOK},
  {"other organization",session.Identity{UserID:"user-a"},nil,map[string]authz.Role{orgA:authz.Owner},orgB,true,http.StatusForbidden},
  {"revoked session",session.Identity{},session.ErrNotFound,map[string]authz.Role{orgA:authz.Owner},orgA,true,http.StatusUnauthorized},
  {"expired session",session.Identity{},session.ErrNotFound,map[string]authz.Role{orgA:authz.Owner},orgA,true,http.StatusUnauthorized},
  {"anonymous",session.Identity{},nil,map[string]authz.Role{orgA:authz.Owner},orgA,false,http.StatusUnauthorized},
 }
 for _,tc:=range tests{
  t.Run(tc.name,func(t *testing.T){
   a:=authz.Authorizer{Sessions:testSession{identity:tc.identity,err:tc.sessionErr},Memberships:testMembership{roles:tc.roles}}
   reached:=false
   protected:=Middleware(a,authz.ReadResponses,func(r *http.Request)string{return r.PathValue("organizationID")},http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
    reached=true
    principal,ok:=Principal(r.Context())
    if !ok||principal.OrganizationID!=orgA||principal.UserID!="user-a"{t.Fatal("invalid principal")}
    w.WriteHeader(http.StatusOK)
   }))
   req:=httptest.NewRequest(http.MethodGet,"/api/organizations/"+tc.requestedOrg+"/workspaces",nil)
   req.SetPathValue("organizationID",tc.requestedOrg)
   if tc.cookie{req.AddCookie(&http.Cookie{Name:session.CookieName,Value:token})}
   rec:=httptest.NewRecorder();protected.ServeHTTP(rec,req)
   if rec.Code!=tc.want{t.Fatalf("status %d, want %d",rec.Code,tc.want)}
   if reached!=(tc.want==http.StatusOK){t.Fatalf("handler reached=%v for status %d",reached,rec.Code)}
  })
 }
}
