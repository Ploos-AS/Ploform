package main

import (
 "html/template"
 "net/http"

 "github.com/Ploos-AS/Ploform/backend/internal/authz"
 "github.com/Ploos-AS/Ploform/backend/internal/httpauth"
)

var inboxPage = template.Must(template.New("inbox").Parse(`<!doctype html>
<html lang="no"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="robots" content="noindex"><title>Ploform – Innboks</title>
<style>
:root{font-family:system-ui,sans-serif;color-scheme:light dark}body{max-width:65rem;margin:2rem auto;padding:0 1rem;line-height:1.5}
header{display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
section{border:1px solid #8887;border-radius:.5rem;padding:1rem;margin-top:1rem}
button{font:inherit;cursor:pointer;padding:.45rem .8rem}
#messages{list-style:none;padding:0;margin:0}#messages li{border-bottom:1px solid #8885}
#messages button{width:100%;text-align:left;border:0;background:transparent;padding:.8rem}
#messages button:hover,#messages button:focus-visible{outline:2px solid currentColor}
small{opacity:.8}pre{white-space:pre-wrap;overflow-wrap:anywhere;font:inherit}
nav{display:flex;align-items:center;gap:1rem;margin-top:1rem}
</style></head><body><main><header><h1>Innboks</h1><small>Skjema: {{.FormID}}</small></header>
<p id="status" role="status">Henter henvendelser …</p>
<section aria-label="Mottatte henvendelser"><ul id="messages"></ul>
<nav><button id="previous" type="button" disabled>Forrige</button><span id="page"></span><button id="next" type="button" disabled>Neste</button></nav></section>
<section aria-label="Valgt henvendelse"><h2>Henvendelse</h2><div id="detail">Velg en henvendelse fra listen.</div></section>
</main>
<script>
"use strict";
const endpoint="/api/organizations/{{.OrganizationID}}/forms/{{.FormID}}/submissions";
const list=document.getElementById("messages"),detail=document.getElementById("detail");
const status=document.getElementById("status"),previous=document.getElementById("previous"),next=document.getElementById("next"),page=document.getElementById("page");
const size=25;let offset=0,items=[];
function text(tag,value){const el=document.createElement(tag);el.textContent=String(value??"");return el;}
function show(item){
 detail.replaceChildren();
 detail.append(text("p",new Date(item.created_at).toLocaleString("nb-NO")));
 for(const [key,value] of Object.entries(item.response||{})){
  const label=text("strong",key+": ");const p=document.createElement("p");
  p.append(label,text("span",Array.isArray(value)?value.join(", "):typeof value==="object"?JSON.stringify(value):value));
  detail.append(p);
 }
}
async function load(){
 status.textContent="Henter henvendelser …";list.replaceChildren();detail.textContent="Velg en henvendelse fra listen.";
 previous.disabled=true;next.disabled=true;
 try{
  const response=await fetch(endpoint+"?limit="+size+"&offset="+offset,{credentials:"same-origin",headers:{"Accept":"application/json"}});
  if(!response.ok)throw Error(response.status===401||response.status===403?"Du har ikke tilgang til innboksen.":"Kunne ikke hente henvendelser ("+response.status+").");
  const data=await response.json();items=data.submissions||[];
  for(const item of items){
   const li=document.createElement("li"),button=document.createElement("button");
   const fields=item.response||{};
   button.append(text("strong",fields.name||fields.email||"Ukjent avsender"),text("span"," — "+(fields.subject||"Uten emne")),document.createElement("br"),text("small",new Date(item.created_at).toLocaleString("nb-NO")));
   button.type="button";button.addEventListener("click",()=>show(item));li.append(button);list.append(li);
  }
  page.textContent="Side "+(Math.floor(offset/size)+1);
  previous.disabled=offset===0;next.disabled=items.length<size;
  status.textContent=items.length?items.length+" henvendelser på denne siden.":"Ingen henvendelser på denne siden.";
 }catch(error){status.textContent=error.message;page.textContent="";}
}
previous.addEventListener("click",()=>{offset=Math.max(0,offset-size);load();});
next.addEventListener("click",()=>{offset+=size;load();});
load();
</script></body></html>`))

func inboxPageRoute(a authz.Authorizer) http.Handler {
 return httpauth.Middleware(a,authz.ReadResponses,func(r *http.Request)string{return r.PathValue("organizationID")},http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  w.Header().Set("Cache-Control","no-store")
  w.Header().Set("X-Content-Type-Options","nosniff")
  w.Header().Set("Content-Security-Policy","default-src 'none'; script-src 'unsafe-inline'; style-src 'unsafe-inline'; connect-src 'self'; base-uri 'none'; frame-ancestors 'none'")
  _=inboxPage.Execute(w,struct{OrganizationID,FormID string}{r.PathValue("organizationID"),r.PathValue("formID")})
 }))
}
