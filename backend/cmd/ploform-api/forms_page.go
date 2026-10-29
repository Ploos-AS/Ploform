package main

import (
 "html/template"
 "net/http"
 "github.com/Ploos-AS/Ploform/backend/internal/authz"
 "github.com/Ploos-AS/Ploform/backend/internal/httpauth"
)

var formsPage = template.Must(template.New("forms").Parse(`<!doctype html>
<html lang="no"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="robots" content="noindex"><title>Ploform – Mine skjemaer</title>
<style>
:root{font-family:system-ui,sans-serif;color-scheme:light dark}body{max-width:60rem;margin:2rem auto;padding:0 1rem;line-height:1.5}
section{border:1px solid #8887;border-radius:.5rem;padding:1rem;margin:1rem 0}
label{display:block;margin:.8rem 0 .3rem}input,select,button{font:inherit;padding:.5rem}
input{box-sizing:border-box;max-width:100%;width:25rem}button{cursor:pointer}
li{margin:.8rem 0;padding:.7rem;border-bottom:1px solid #8885}ul{list-style:none;padding:0}
a{margin-right:.7rem}small{opacity:.8}.actions{display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
</style></head><body><main><h1>Mine skjemaer</h1><p id="status" role="status">Laster arbeidsområder …</p>
<section><label for="workspace">Arbeidsområde</label><select id="workspace"></select></section>
<section><h2>Nytt kontaktskjema</h2><form id="create"><label for="title">Tittel</label><input id="title" required maxlength="120" placeholder="Kontakt oss"><label for="slug">Kort navn (URL)</label><input id="slug" required pattern="[a-z0-9][a-z0-9-]*" maxlength="63" placeholder="kontakt"><p><button type="submit">Opprett kontaktskjema</button></p></form></section>
<section><h2>Skjemaer</h2><ul id="forms"></ul></section></main>
<script>
"use strict";
const org="{{.OrganizationID}}",base="/api/organizations/"+org;
const workspace=document.getElementById("workspace"),list=document.getElementById("forms"),status=document.getElementById("status");
function el(tag,value){const node=document.createElement(tag);node.textContent=String(value??"");return node;}
async function get(url){const r=await fetch(url,{credentials:"same-origin"});if(!r.ok)throw Error("Forespørselen feilet ("+r.status+")");return r.json();}
async function write(url,data){
 const csrf=await get("/api/csrf");
 const r=await fetch(url,{method:"POST",credentials:"same-origin",headers:{"Content-Type":"application/json","X-CSRF-Token":csrf.csrf_token},body:JSON.stringify(data)});
 if(!r.ok)throw Error("Kunne ikke lagre ("+r.status+")");
 return r.json();
}
async function loadForms(){
 list.replaceChildren();const ws=workspace.value;if(!ws)return;
 try{
  const data=await get(base+"/workspaces/"+ws+"/forms");
  for(const form of data.forms||[]){
   const li=el("li",""),title=el("strong",form.title+" "),state=el("small","("+form.status+")"),actions=el("div","");
   actions.className="actions";li.append(title,state,actions);
   const inbox=el("a","Åpne innboks");inbox.href="/inbox/"+encodeURIComponent(org)+"/"+encodeURIComponent(form.id);actions.append(inbox);
   if(form.status==="draft"){
    const publish=el("button","Publiser");publish.type="button";
    publish.addEventListener("click",async()=>{try{await write(base+"/forms/"+form.id+"/publish",{});status.textContent="Skjema publisert.";await loadForms();}catch(e){status.textContent=e.message;}});
    actions.append(publish);
   }
   if(form.status==="published"){
    const link=el("a","Åpne offentlig skjema");link.href="/f/"+encodeURIComponent(org)+"/"+encodeURIComponent(workspace.selectedOptions[0].dataset.slug)+"/"+encodeURIComponent(form.slug);link.target="_blank";link.rel="noopener noreferrer";actions.append(link);
   }
   list.append(li);
  }
  status.textContent=(data.forms||[]).length+" skjemaer.";
 }catch(e){status.textContent=e.message;}
}
async function start(){
 try{
  const data=await get(base+"/workspaces");
  for(const ws of data.workspaces||[]){const opt=el("option",ws.name||ws.slug);opt.value=ws.id;opt.dataset.slug=ws.slug;workspace.append(opt);}
  if(!workspace.options.length){status.textContent="Ingen arbeidsområder tilgjengelig.";return;}
  await loadForms();
 }catch(e){status.textContent=e.message;}
}
workspace.addEventListener("change",loadForms);
document.getElementById("create").addEventListener("submit",async e=>{
 e.preventDefault();try{
  await write(base+"/workspaces/"+workspace.value+"/contact-forms",{title:document.getElementById("title").value,slug:document.getElementById("slug").value});
  e.target.reset();await loadForms();status.textContent="Kontaktskjema opprettet.";
 }catch(error){status.textContent=error.message;}
});
start();
</script></body></html>`))

func formsPageRoute(a authz.Authorizer) http.Handler {
 return httpauth.Middleware(a,authz.ReadResponses,func(r *http.Request)string{return r.PathValue("organizationID")},http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  w.Header().Set("Cache-Control","no-store")
  w.Header().Set("X-Content-Type-Options","nosniff")
  w.Header().Set("Content-Security-Policy","default-src 'none'; script-src 'unsafe-inline'; style-src 'unsafe-inline'; connect-src 'self'; base-uri 'none'; frame-ancestors 'none'")
  _=formsPage.Execute(w,struct{OrganizationID string}{r.PathValue("organizationID")})
 }))
}
