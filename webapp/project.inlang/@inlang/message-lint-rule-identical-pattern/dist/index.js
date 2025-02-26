var an=Object.create;var Je=Object.defineProperty;var Tn=Object.getOwnPropertyDescriptor;var cn=Object.getOwnPropertyNames;var ln=Object.getPrototypeOf,pn=Object.prototype.hasOwnProperty;var dn=(c,t)=>()=>(t||c((t={exports:{}}).exports,t),t.exports);var yn=(c,t,u,T)=>{if(t&&typeof t=="object"||typeof t=="function")for(let p of cn(t))!pn.call(c,p)&&p!==u&&Je(c,p,{get:()=>t[p],enumerable:!(T=Tn(t,p))||T.enumerable});return c};var fn=(c,t,u)=>(u=c!=null?an(ln(c)):{},yn(t||!c||!c.__esModule?Je(u,"default",{value:c,enumerable:!0}):u,c));var Ye=dn(i=>{"use strict";Object.defineProperty(i,"__esModule",{value:!0});i.Type=i.JsonType=i.JavaScriptTypeBuilder=i.JsonTypeBuilder=i.TypeBuilder=i.TypeBuilderError=i.TransformEncodeBuilder=i.TransformDecodeBuilder=i.TemplateLiteralDslParser=i.TemplateLiteralGenerator=i.TemplateLiteralGeneratorError=i.TemplateLiteralFinite=i.TemplateLiteralFiniteError=i.TemplateLiteralParser=i.TemplateLiteralParserError=i.TemplateLiteralResolver=i.TemplateLiteralPattern=i.TemplateLiteralPatternError=i.UnionResolver=i.KeyArrayResolver=i.KeyArrayResolverError=i.KeyResolver=i.ObjectMap=i.Intrinsic=i.IndexedAccessor=i.TypeClone=i.TypeExtends=i.TypeExtendsResult=i.TypeExtendsError=i.ExtendsUndefined=i.TypeGuard=i.TypeGuardUnknownTypeError=i.ValueGuard=i.FormatRegistry=i.TypeBoxError=i.TypeRegistry=i.PatternStringExact=i.PatternNumberExact=i.PatternBooleanExact=i.PatternString=i.PatternNumber=i.PatternBoolean=i.Kind=i.Hint=i.Optional=i.Readonly=i.Transform=void 0;i.Transform=Symbol.for("TypeBox.Transform");i.Readonly=Symbol.for("TypeBox.Readonly");i.Optional=Symbol.for("TypeBox.Optional");i.Hint=Symbol.for("TypeBox.Hint");i.Kind=Symbol.for("TypeBox.Kind");i.PatternBoolean="(true|false)";i.PatternNumber="(0|[1-9][0-9]*)";i.PatternString="(.*)";i.PatternBooleanExact=`^${i.PatternBoolean}$`;i.PatternNumberExact=`^${i.PatternNumber}$`;i.PatternStringExact=`^${i.PatternString}$`;var ke;(function(c){let t=new Map;function u(){return new Map(t)}c.Entries=u;function T(){return t.clear()}c.Clear=T;function p(I){return t.delete(I)}c.Delete=p;function a(I){return t.has(I)}c.Has=a;function s(I,R){t.set(I,R)}c.Set=s;function y(I){return t.get(I)}c.Get=y})(ke||(i.TypeRegistry=ke={}));var k=class extends Error{constructor(t){super(t)}};i.TypeBoxError=k;var Qe;(function(c){let t=new Map;function u(){return new Map(t)}c.Entries=u;function T(){return t.clear()}c.Clear=T;function p(I){return t.delete(I)}c.Delete=p;function a(I){return t.has(I)}c.Has=a;function s(I,R){t.set(I,R)}c.Set=s;function y(I){return t.get(I)}c.Get=y})(Qe||(i.FormatRegistry=Qe={}));var f;(function(c){function t(O){return Array.isArray(O)}c.IsArray=t;function u(O){return typeof O=="bigint"}c.IsBigInt=u;function T(O){return typeof O=="boolean"}c.IsBoolean=T;function p(O){return O instanceof globalThis.Date}c.IsDate=p;function a(O){return O===null}c.IsNull=a;function s(O){return typeof O=="number"}c.IsNumber=s;function y(O){return typeof O=="object"&&O!==null}c.IsObject=y;function I(O){return typeof O=="string"}c.IsString=I;function R(O){return O instanceof globalThis.Uint8Array}c.IsUint8Array=R;function b(O){return O===void 0}c.IsUndefined=b})(f||(i.ValueGuard=f={}));var Me=class extends k{};i.TypeGuardUnknownTypeError=Me;var o;(function(c){function t(r){try{return new RegExp(r),!0}catch{return!1}}function u(r){if(!f.IsString(r))return!1;for(let C=0;C<r.length;C++){let B=r.charCodeAt(C);if(B>=7&&B<=13||B===27||B===127)return!1}return!0}function T(r){return s(r)||$(r)}function p(r){return f.IsUndefined(r)||f.IsBigInt(r)}function a(r){return f.IsUndefined(r)||f.IsNumber(r)}function s(r){return f.IsUndefined(r)||f.IsBoolean(r)}function y(r){return f.IsUndefined(r)||f.IsString(r)}function I(r){return f.IsUndefined(r)||f.IsString(r)&&u(r)&&t(r)}function R(r){return f.IsUndefined(r)||f.IsString(r)&&u(r)}function b(r){return f.IsUndefined(r)||$(r)}function O(r){return L(r,"Any")&&y(r.$id)}c.TAny=O;function m(r){return L(r,"Array")&&r.type==="array"&&y(r.$id)&&$(r.items)&&a(r.minItems)&&a(r.maxItems)&&s(r.uniqueItems)&&b(r.contains)&&a(r.minContains)&&a(r.maxContains)}c.TArray=m;function d(r){return L(r,"AsyncIterator")&&r.type==="AsyncIterator"&&y(r.$id)&&$(r.items)}c.TAsyncIterator=d;function U(r){return L(r,"BigInt")&&r.type==="bigint"&&y(r.$id)&&p(r.exclusiveMaximum)&&p(r.exclusiveMinimum)&&p(r.maximum)&&p(r.minimum)&&p(r.multipleOf)}c.TBigInt=U;function N(r){return L(r,"Boolean")&&r.type==="boolean"&&y(r.$id)}c.TBoolean=N;function S(r){return L(r,"Constructor")&&r.type==="Constructor"&&y(r.$id)&&f.IsArray(r.parameters)&&r.parameters.every(C=>$(C))&&$(r.returns)}c.TConstructor=S;function w(r){return L(r,"Date")&&r.type==="Date"&&y(r.$id)&&a(r.exclusiveMaximumTimestamp)&&a(r.exclusiveMinimumTimestamp)&&a(r.maximumTimestamp)&&a(r.minimumTimestamp)&&a(r.multipleOfTimestamp)}c.TDate=w;function v(r){return L(r,"Function")&&r.type==="Function"&&y(r.$id)&&f.IsArray(r.parameters)&&r.parameters.every(C=>$(C))&&$(r.returns)}c.TFunction=v;function j(r){return L(r,"Integer")&&r.type==="integer"&&y(r.$id)&&a(r.exclusiveMaximum)&&a(r.exclusiveMinimum)&&a(r.maximum)&&a(r.minimum)&&a(r.multipleOf)}c.TInteger=j;function F(r){return L(r,"Intersect")&&!(f.IsString(r.type)&&r.type!=="object")&&f.IsArray(r.allOf)&&r.allOf.every(C=>$(C)&&!re(C))&&y(r.type)&&(s(r.unevaluatedProperties)||b(r.unevaluatedProperties))&&y(r.$id)}c.TIntersect=F;function Te(r){return L(r,"Iterator")&&r.type==="Iterator"&&y(r.$id)&&$(r.items)}c.TIterator=Te;function L(r,C){return G(r)&&r[i.Kind]===C}c.TKindOf=L;function G(r){return f.IsObject(r)&&i.Kind in r&&f.IsString(r[i.Kind])}c.TKind=G;function h(r){return V(r)&&f.IsString(r.const)}c.TLiteralString=h;function ce(r){return V(r)&&f.IsNumber(r.const)}c.TLiteralNumber=ce;function je(r){return V(r)&&f.IsBoolean(r.const)}c.TLiteralBoolean=je;function V(r){return L(r,"Literal")&&y(r.$id)&&(f.IsBoolean(r.const)||f.IsNumber(r.const)||f.IsString(r.const))}c.TLiteral=V;function le(r){return L(r,"Never")&&f.IsObject(r.not)&&Object.getOwnPropertyNames(r.not).length===0}c.TNever=le;function K(r){return L(r,"Not")&&$(r.not)}c.TNot=K;function ee(r){return L(r,"Null")&&r.type==="null"&&y(r.$id)}c.TNull=ee;function ne(r){return L(r,"Number")&&r.type==="number"&&y(r.$id)&&a(r.exclusiveMaximum)&&a(r.exclusiveMinimum)&&a(r.maximum)&&a(r.minimum)&&a(r.multipleOf)}c.TNumber=ne;function J(r){return L(r,"Object")&&r.type==="object"&&y(r.$id)&&f.IsObject(r.properties)&&T(r.additionalProperties)&&a(r.minProperties)&&a(r.maxProperties)&&Object.entries(r.properties).every(([C,B])=>u(C)&&$(B))}c.TObject=J;function te(r){return L(r,"Promise")&&r.type==="Promise"&&y(r.$id)&&$(r.item)}c.TPromise=te;function pe(r){return L(r,"Record")&&r.type==="object"&&y(r.$id)&&T(r.additionalProperties)&&f.IsObject(r.patternProperties)&&(C=>{let B=Object.getOwnPropertyNames(C.patternProperties);return B.length===1&&t(B[0])&&f.IsObject(C.patternProperties)&&$(C.patternProperties[B[0]])})(r)}c.TRecord=pe;function ge(r){return f.IsObject(r)&&i.Hint in r&&r[i.Hint]==="Recursive"}c.TRecursive=ge;function de(r){return L(r,"Ref")&&y(r.$id)&&f.IsString(r.$ref)}c.TRef=de;function ye(r){return L(r,"String")&&r.type==="string"&&y(r.$id)&&a(r.minLength)&&a(r.maxLength)&&I(r.pattern)&&R(r.format)}c.TString=ye;function fe(r){return L(r,"Symbol")&&r.type==="symbol"&&y(r.$id)}c.TSymbol=fe;function q(r){return L(r,"TemplateLiteral")&&r.type==="string"&&f.IsString(r.pattern)&&r.pattern[0]==="^"&&r.pattern[r.pattern.length-1]==="$"}c.TTemplateLiteral=q;function Ie(r){return L(r,"This")&&y(r.$id)&&f.IsString(r.$ref)}c.TThis=Ie;function re(r){return f.IsObject(r)&&i.Transform in r}c.TTransform=re;function g(r){return L(r,"Tuple")&&r.type==="array"&&y(r.$id)&&f.IsNumber(r.minItems)&&f.IsNumber(r.maxItems)&&r.minItems===r.maxItems&&(f.IsUndefined(r.items)&&f.IsUndefined(r.additionalItems)&&r.minItems===0||f.IsArray(r.items)&&r.items.every(C=>$(C)))}c.TTuple=g;function Oe(r){return L(r,"Undefined")&&r.type==="undefined"&&y(r.$id)}c.TUndefined=Oe;function $e(r){return H(r)&&r.anyOf.every(C=>h(C)||ce(C))}c.TUnionLiteral=$e;function H(r){return L(r,"Union")&&y(r.$id)&&f.IsObject(r)&&f.IsArray(r.anyOf)&&r.anyOf.every(C=>$(C))}c.TUnion=H;function _(r){return L(r,"Uint8Array")&&r.type==="Uint8Array"&&y(r.$id)&&a(r.minByteLength)&&a(r.maxByteLength)}c.TUint8Array=_;function E(r){return L(r,"Unknown")&&y(r.$id)}c.TUnknown=E;function be(r){return L(r,"Unsafe")}c.TUnsafe=be;function ie(r){return L(r,"Void")&&r.type==="void"&&y(r.$id)}c.TVoid=ie;function Ke(r){return f.IsObject(r)&&r[i.Readonly]==="Readonly"}c.TReadonly=Ke;function Fe(r){return f.IsObject(r)&&r[i.Optional]==="Optional"}c.TOptional=Fe;function $(r){return f.IsObject(r)&&(O(r)||m(r)||N(r)||U(r)||d(r)||S(r)||w(r)||v(r)||j(r)||F(r)||Te(r)||V(r)||le(r)||K(r)||ee(r)||ne(r)||J(r)||te(r)||pe(r)||de(r)||ye(r)||fe(r)||q(r)||Ie(r)||g(r)||Oe(r)||H(r)||_(r)||E(r)||be(r)||ie(r)||G(r)&&ke.Has(r[i.Kind]))}c.TSchema=$})(o||(i.TypeGuard=o={}));var Xe;(function(c){function t(u){return u[i.Kind]==="Intersect"?u.allOf.every(T=>t(T)):u[i.Kind]==="Union"?u.anyOf.some(T=>t(T)):u[i.Kind]==="Undefined"?!0:u[i.Kind]==="Not"?!t(u.not):!1}c.Check=t})(Xe||(i.ExtendsUndefined=Xe={}));var Ue=class extends k{};i.TypeExtendsError=Ue;var l;(function(c){c[c.Union=0]="Union",c[c.True=1]="True",c[c.False=2]="False"})(l||(i.TypeExtendsResult=l={}));var z;(function(c){function t(e){return e===l.False?e:l.True}function u(e){throw new Ue(e)}function T(e){return o.TNever(e)||o.TIntersect(e)||o.TUnion(e)||o.TUnknown(e)||o.TAny(e)}function p(e,n){return o.TNever(n)?L(e,n):o.TIntersect(n)?v(e,n):o.TUnion(n)?Ee(e,n):o.TUnknown(n)?ze(e,n):o.TAny(n)?a(e,n):u("StructuralRight")}function a(e,n){return l.True}function s(e,n){return o.TIntersect(n)?v(e,n):o.TUnion(n)&&n.anyOf.some(A=>o.TAny(A)||o.TUnknown(A))?l.True:o.TUnion(n)?l.Union:o.TUnknown(n)||o.TAny(n)?l.True:l.Union}function y(e,n){return o.TUnknown(e)?l.False:o.TAny(e)?l.Union:o.TNever(e)?l.True:l.False}function I(e,n){return o.TObject(n)&&q(n)?l.True:T(n)?p(e,n):o.TArray(n)?t(x(e.items,n.items)):l.False}function R(e,n){return T(n)?p(e,n):o.TAsyncIterator(n)?t(x(e.items,n.items)):l.False}function b(e,n){return T(n)?p(e,n):o.TObject(n)?g(e,n):o.TRecord(n)?E(e,n):o.TBigInt(n)?l.True:l.False}function O(e,n){return o.TLiteral(e)&&f.IsBoolean(e.const)||o.TBoolean(e)?l.True:l.False}function m(e,n){return T(n)?p(e,n):o.TObject(n)?g(e,n):o.TRecord(n)?E(e,n):o.TBoolean(n)?l.True:l.False}function d(e,n){return T(n)?p(e,n):o.TObject(n)?g(e,n):o.TConstructor(n)?e.parameters.length>n.parameters.length?l.False:e.parameters.every((A,D)=>t(x(n.parameters[D],A))===l.True)?t(x(e.returns,n.returns)):l.False:l.False}function U(e,n){return T(n)?p(e,n):o.TObject(n)?g(e,n):o.TRecord(n)?E(e,n):o.TDate(n)?l.True:l.False}function N(e,n){return T(n)?p(e,n):o.TObject(n)?g(e,n):o.TFunction(n)?e.parameters.length>n.parameters.length?l.False:e.parameters.every((A,D)=>t(x(n.parameters[D],A))===l.True)?t(x(e.returns,n.returns)):l.False:l.False}function S(e,n){return o.TLiteral(e)&&f.IsNumber(e.const)||o.TNumber(e)||o.TInteger(e)?l.True:l.False}function w(e,n){return o.TInteger(n)||o.TNumber(n)?l.True:T(n)?p(e,n):o.TObject(n)?g(e,n):o.TRecord(n)?E(e,n):l.False}function v(e,n){return n.allOf.every(A=>x(e,A)===l.True)?l.True:l.False}function j(e,n){return e.allOf.some(A=>x(A,n)===l.True)?l.True:l.False}function F(e,n){return T(n)?p(e,n):o.TIterator(n)?t(x(e.items,n.items)):l.False}function Te(e,n){return o.TLiteral(n)&&n.const===e.const?l.True:T(n)?p(e,n):o.TObject(n)?g(e,n):o.TRecord(n)?E(e,n):o.TString(n)?ie(e,n):o.TNumber(n)?V(e,n):o.TInteger(n)?S(e,n):o.TBoolean(n)?O(e,n):l.False}function L(e,n){return l.False}function G(e,n){return l.True}function h(e){let[n,A]=[e,0];for(;o.TNot(n);)n=n.not,A+=1;return A%2===0?n:i.Type.Unknown()}function ce(e,n){return o.TNot(e)?x(h(e),n):o.TNot(n)?x(e,h(n)):u("Invalid fallthrough for Not")}function je(e,n){return T(n)?p(e,n):o.TObject(n)?g(e,n):o.TRecord(n)?E(e,n):o.TNull(n)?l.True:l.False}function V(e,n){return o.TLiteralNumber(e)||o.TNumber(e)||o.TInteger(e)?l.True:l.False}function le(e,n){return T(n)?p(e,n):o.TObject(n)?g(e,n):o.TRecord(n)?E(e,n):o.TInteger(n)||o.TNumber(n)?l.True:l.False}function K(e,n){return Object.getOwnPropertyNames(e.properties).length===n}function ee(e){return q(e)}function ne(e){return K(e,0)||K(e,1)&&"description"in e.properties&&o.TUnion(e.properties.description)&&e.properties.description.anyOf.length===2&&(o.TString(e.properties.description.anyOf[0])&&o.TUndefined(e.properties.description.anyOf[1])||o.TString(e.properties.description.anyOf[1])&&o.TUndefined(e.properties.description.anyOf[0]))}function J(e){return K(e,0)}function te(e){return K(e,0)}function pe(e){return K(e,0)}function ge(e){return K(e,0)}function de(e){return q(e)}function ye(e){let n=i.Type.Number();return K(e,0)||K(e,1)&&"length"in e.properties&&t(x(e.properties.length,n))===l.True}function fe(e){return K(e,0)}function q(e){let n=i.Type.Number();return K(e,0)||K(e,1)&&"length"in e.properties&&t(x(e.properties.length,n))===l.True}function Ie(e){let n=i.Type.Function([i.Type.Any()],i.Type.Any());return K(e,0)||K(e,1)&&"then"in e.properties&&t(x(e.properties.then,n))===l.True}function re(e,n){return x(e,n)===l.False||o.TOptional(e)&&!o.TOptional(n)?l.False:l.True}function g(e,n){return o.TUnknown(e)?l.False:o.TAny(e)?l.Union:o.TNever(e)||o.TLiteralString(e)&&ee(n)||o.TLiteralNumber(e)&&J(n)||o.TLiteralBoolean(e)&&te(n)||o.TSymbol(e)&&ne(n)||o.TBigInt(e)&&pe(n)||o.TString(e)&&ee(n)||o.TSymbol(e)&&ne(n)||o.TNumber(e)&&J(n)||o.TInteger(e)&&J(n)||o.TBoolean(e)&&te(n)||o.TUint8Array(e)&&de(n)||o.TDate(e)&&ge(n)||o.TConstructor(e)&&fe(n)||o.TFunction(e)&&ye(n)?l.True:o.TRecord(e)&&o.TString(H(e))?n[i.Hint]==="Record"?l.True:l.False:o.TRecord(e)&&o.TNumber(H(e))?K(n,0)?l.True:l.False:l.False}function Oe(e,n){return T(n)?p(e,n):o.TRecord(n)?E(e,n):o.TObject(n)?(()=>{for(let A of Object.getOwnPropertyNames(n.properties)){if(!(A in e.properties)&&!o.TOptional(n.properties[A]))return l.False;if(o.TOptional(n.properties[A]))return l.True;if(re(e.properties[A],n.properties[A])===l.False)return l.False}return l.True})():l.False}function $e(e,n){return T(n)?p(e,n):o.TObject(n)&&Ie(n)?l.True:o.TPromise(n)?t(x(e.item,n.item)):l.False}function H(e){return i.PatternNumberExact in e.patternProperties?i.Type.Number():i.PatternStringExact in e.patternProperties?i.Type.String():u("Unknown record key pattern")}function _(e){return i.PatternNumberExact in e.patternProperties?e.patternProperties[i.PatternNumberExact]:i.PatternStringExact in e.patternProperties?e.patternProperties[i.PatternStringExact]:u("Unable to get record value schema")}function E(e,n){let[A,D]=[H(n),_(n)];return o.TLiteralString(e)&&o.TNumber(A)&&t(x(e,D))===l.True?l.True:o.TUint8Array(e)&&o.TNumber(A)||o.TString(e)&&o.TNumber(A)||o.TArray(e)&&o.TNumber(A)?x(e,D):o.TObject(e)?(()=>{for(let sn of Object.getOwnPropertyNames(e.properties))if(re(D,e.properties[sn])===l.False)return l.False;return l.True})():l.False}function be(e,n){return T(n)?p(e,n):o.TObject(n)?g(e,n):o.TRecord(n)?x(_(e),_(n)):l.False}function ie(e,n){return o.TLiteral(e)&&f.IsString(e.const)||o.TString(e)?l.True:l.False}function Ke(e,n){return T(n)?p(e,n):o.TObject(n)?g(e,n):o.TRecord(n)?E(e,n):o.TString(n)?l.True:l.False}function Fe(e,n){return T(n)?p(e,n):o.TObject(n)?g(e,n):o.TRecord(n)?E(e,n):o.TSymbol(n)?l.True:l.False}function $(e,n){return o.TTemplateLiteral(e)?x(M.Resolve(e),n):o.TTemplateLiteral(n)?x(e,M.Resolve(n)):u("Invalid fallthrough for TemplateLiteral")}function r(e,n){return o.TArray(n)&&e.items!==void 0&&e.items.every(A=>x(A,n.items)===l.True)}function C(e,n){return o.TNever(e)?l.True:o.TUnknown(e)?l.False:o.TAny(e)?l.Union:l.False}function B(e,n){return T(n)?p(e,n):o.TObject(n)&&q(n)||o.TArray(n)&&r(e,n)?l.True:o.TTuple(n)?f.IsUndefined(e.items)&&!f.IsUndefined(n.items)||!f.IsUndefined(e.items)&&f.IsUndefined(n.items)?l.False:f.IsUndefined(e.items)&&!f.IsUndefined(n.items)||e.items.every((A,D)=>x(A,n.items[D])===l.True)?l.True:l.False:l.False}function he(e,n){return T(n)?p(e,n):o.TObject(n)?g(e,n):o.TRecord(n)?E(e,n):o.TUint8Array(n)?l.True:l.False}function en(e,n){return T(n)?p(e,n):o.TObject(n)?g(e,n):o.TRecord(n)?E(e,n):o.TVoid(n)?rn(e,n):o.TUndefined(n)?l.True:l.False}function Ee(e,n){return n.anyOf.some(A=>x(e,A)===l.True)?l.True:l.False}function nn(e,n){return e.anyOf.every(A=>x(A,n)===l.True)?l.True:l.False}function ze(e,n){return l.True}function tn(e,n){return o.TNever(n)?L(e,n):o.TIntersect(n)?v(e,n):o.TUnion(n)?Ee(e,n):o.TAny(n)?a(e,n):o.TString(n)?ie(e,n):o.TNumber(n)?V(e,n):o.TInteger(n)?S(e,n):o.TBoolean(n)?O(e,n):o.TArray(n)?y(e,n):o.TTuple(n)?C(e,n):o.TObject(n)?g(e,n):o.TUnknown(n)?l.True:l.False}function rn(e,n){return o.TUndefined(e)||o.TUndefined(e)?l.True:l.False}function on(e,n){return o.TIntersect(n)?v(e,n):o.TUnion(n)?Ee(e,n):o.TUnknown(n)?ze(e,n):o.TAny(n)?a(e,n):o.TObject(n)?g(e,n):o.TVoid(n)?l.True:l.False}function x(e,n){return o.TTemplateLiteral(e)||o.TTemplateLiteral(n)?$(e,n):o.TNot(e)||o.TNot(n)?ce(e,n):o.TAny(e)?s(e,n):o.TArray(e)?I(e,n):o.TBigInt(e)?b(e,n):o.TBoolean(e)?m(e,n):o.TAsyncIterator(e)?R(e,n):o.TConstructor(e)?d(e,n):o.TDate(e)?U(e,n):o.TFunction(e)?N(e,n):o.TInteger(e)?w(e,n):o.TIntersect(e)?j(e,n):o.TIterator(e)?F(e,n):o.TLiteral(e)?Te(e,n):o.TNever(e)?G(e,n):o.TNull(e)?je(e,n):o.TNumber(e)?le(e,n):o.TObject(e)?Oe(e,n):o.TRecord(e)?be(e,n):o.TString(e)?Ke(e,n):o.TSymbol(e)?Fe(e,n):o.TTuple(e)?B(e,n):o.TPromise(e)?$e(e,n):o.TUint8Array(e)?he(e,n):o.TUndefined(e)?en(e,n):o.TUnion(e)?nn(e,n):o.TUnknown(e)?tn(e,n):o.TVoid(e)?on(e,n):u(`Unknown left type operand '${e[i.Kind]}'`)}function un(e,n){return x(e,n)}c.Extends=un})(z||(i.TypeExtends=z={}));var P;(function(c){function t(I){return I.map(R=>a(R))}function u(I){return new Date(I.getTime())}function T(I){return new Uint8Array(I)}function p(I){let R=Object.getOwnPropertyNames(I).reduce((O,m)=>({...O,[m]:a(I[m])}),{}),b=Object.getOwnPropertySymbols(I).reduce((O,m)=>({...O,[m]:a(I[m])}),{});return{...R,...b}}function a(I){return f.IsArray(I)?t(I):f.IsDate(I)?u(I):f.IsUint8Array(I)?T(I):f.IsObject(I)?p(I):I}function s(I){return I.map(R=>y(R))}c.Rest=s;function y(I,R={}){return{...a(I),...R}}c.Type=y})(P||(i.TypeClone=P={}));var Ve;(function(c){function t(d){return d.map(U=>{let{[i.Optional]:N,...S}=P.Type(U);return S})}function u(d){return d.every(U=>o.TOptional(U))}function T(d){return d.some(U=>o.TOptional(U))}function p(d){return u(d.allOf)?i.Type.Optional(i.Type.Intersect(t(d.allOf))):d}function a(d){return T(d.anyOf)?i.Type.Optional(i.Type.Union(t(d.anyOf))):d}function s(d){return d[i.Kind]==="Intersect"?p(d):d[i.Kind]==="Union"?a(d):d}function y(d,U){let N=d.allOf.reduce((S,w)=>{let v=O(w,U);return v[i.Kind]==="Never"?S:[...S,v]},[]);return s(i.Type.Intersect(N))}function I(d,U){let N=d.anyOf.map(S=>O(S,U));return s(i.Type.Union(N))}function R(d,U){let N=d.properties[U];return f.IsUndefined(N)?i.Type.Never():i.Type.Union([N])}function b(d,U){let N=d.items;if(f.IsUndefined(N))return i.Type.Never();let S=N[U];return f.IsUndefined(S)?i.Type.Never():S}function O(d,U){return d[i.Kind]==="Intersect"?y(d,U):d[i.Kind]==="Union"?I(d,U):d[i.Kind]==="Object"?R(d,U):d[i.Kind]==="Tuple"?b(d,U):i.Type.Never()}function m(d,U,N={}){let S=U.map(w=>O(d,w.toString()));return s(i.Type.Union(S,N))}c.Resolve=m})(Ve||(i.IndexedAccessor=Ve={}));var W;(function(c){function t(b){let[O,m]=[b.slice(0,1),b.slice(1)];return`${O.toLowerCase()}${m}`}function u(b){let[O,m]=[b.slice(0,1),b.slice(1)];return`${O.toUpperCase()}${m}`}function T(b){return b.toUpperCase()}function p(b){return b.toLowerCase()}function a(b,O){let m=X.ParseExact(b.pattern);if(!Y.Check(m))return{...b,pattern:s(b.pattern,O)};let N=[...Z.Generate(m)].map(v=>i.Type.Literal(v)),S=y(N,O),w=i.Type.Union(S);return i.Type.TemplateLiteral([w])}function s(b,O){return typeof b=="string"?O==="Uncapitalize"?t(b):O==="Capitalize"?u(b):O==="Uppercase"?T(b):O==="Lowercase"?p(b):b:b.toString()}function y(b,O){if(b.length===0)return[];let[m,...d]=b;return[R(m,O),...y(d,O)]}function I(b,O){return o.TTemplateLiteral(b)?a(b,O):o.TUnion(b)?i.Type.Union(y(b.anyOf,O)):o.TLiteral(b)?i.Type.Literal(s(b.const,O)):b}function R(b,O){return I(b,O)}c.Map=R})(W||(i.Intrinsic=W={}));var Q;(function(c){function t(s,y){return i.Type.Intersect(s.allOf.map(I=>p(I,y)),{...s})}function u(s,y){return i.Type.Union(s.anyOf.map(I=>p(I,y)),{...s})}function T(s,y){return y(s)}function p(s,y){return s[i.Kind]==="Intersect"?t(s,y):s[i.Kind]==="Union"?u(s,y):s[i.Kind]==="Object"?T(s,y):s}function a(s,y,I){return{...p(P.Type(s),y),...I}}c.Map=a})(Q||(i.ObjectMap=Q={}));var Re;(function(c){function t(R){return R[0]==="^"&&R[R.length-1]==="$"?R.slice(1,R.length-1):R}function u(R,b){return R.allOf.reduce((O,m)=>[...O,...s(m,b)],[])}function T(R,b){let O=R.anyOf.map(m=>s(m,b));return[...O.reduce((m,d)=>d.map(U=>O.every(N=>N.includes(U))?m.add(U):m)[0],new Set)]}function p(R,b){return Object.getOwnPropertyNames(R.properties)}function a(R,b){return b.includePatterns?Object.getOwnPropertyNames(R.patternProperties):[]}function s(R,b){return o.TIntersect(R)?u(R,b):o.TUnion(R)?T(R,b):o.TObject(R)?p(R,b):o.TRecord(R)?a(R,b):[]}function y(R,b){return[...new Set(s(R,b))]}c.ResolveKeys=y;function I(R){return`^(${y(R,{includePatterns:!0}).map(m=>`(${t(m)})`).join("|")})$`}c.ResolvePattern=I})(Re||(i.KeyResolver=Re={}));var Pe=class extends k{};i.KeyArrayResolverError=Pe;var oe;(function(c){function t(u){return Array.isArray(u)?u:o.TUnionLiteral(u)?u.anyOf.map(T=>T.const.toString()):o.TLiteral(u)?[u.const]:o.TTemplateLiteral(u)?(()=>{let T=X.ParseExact(u.pattern);if(!Y.Check(T))throw new Pe("Cannot resolve keys from infinite template expression");return[...Z.Generate(T)]})():[]}c.Resolve=t})(oe||(i.KeyArrayResolver=oe={}));var qe;(function(c){function*t(T){for(let p of T.anyOf)p[i.Kind]==="Union"?yield*t(p):yield p}function u(T){return i.Type.Union([...t(T)],{...T})}c.Resolve=u})(qe||(i.UnionResolver=qe={}));var Ne=class extends k{};i.TemplateLiteralPatternError=Ne;var Le;(function(c){function t(a){throw new Ne(a)}function u(a){return a.replace(/[.*+?^${}()|[\]\\]/g,"\\$&")}function T(a,s){return o.TTemplateLiteral(a)?a.pattern.slice(1,a.pattern.length-1):o.TUnion(a)?`(${a.anyOf.map(y=>T(y,s)).join("|")})`:o.TNumber(a)?`${s}${i.PatternNumber}`:o.TInteger(a)?`${s}${i.PatternNumber}`:o.TBigInt(a)?`${s}${i.PatternNumber}`:o.TString(a)?`${s}${i.PatternString}`:o.TLiteral(a)?`${s}${u(a.const.toString())}`:o.TBoolean(a)?`${s}${i.PatternBoolean}`:t(`Unexpected Kind '${a[i.Kind]}'`)}function p(a){return`^${a.map(s=>T(s,"")).join("")}$`}c.Create=p})(Le||(i.TemplateLiteralPattern=Le={}));var M;(function(c){function t(u){let T=X.ParseExact(u.pattern);if(!Y.Check(T))return i.Type.String();let p=[...Z.Generate(T)].map(a=>i.Type.Literal(a));return i.Type.Union(p)}c.Resolve=t})(M||(i.TemplateLiteralResolver=M={}));var ue=class extends k{};i.TemplateLiteralParserError=ue;var X;(function(c){function t(d,U,N){return d[U]===N&&d.charCodeAt(U-1)!==92}function u(d,U){return t(d,U,"(")}function T(d,U){return t(d,U,")")}function p(d,U){return t(d,U,"|")}function a(d){if(!(u(d,0)&&T(d,d.length-1)))return!1;let U=0;for(let N=0;N<d.length;N++)if(u(d,N)&&(U+=1),T(d,N)&&(U-=1),U===0&&N!==d.length-1)return!1;return!0}function s(d){return d.slice(1,d.length-1)}function y(d){let U=0;for(let N=0;N<d.length;N++)if(u(d,N)&&(U+=1),T(d,N)&&(U-=1),p(d,N)&&U===0)return!0;return!1}function I(d){for(let U=0;U<d.length;U++)if(u(d,U))return!0;return!1}function R(d){let[U,N]=[0,0],S=[];for(let v=0;v<d.length;v++)if(u(d,v)&&(U+=1),T(d,v)&&(U-=1),p(d,v)&&U===0){let j=d.slice(N,v);j.length>0&&S.push(O(j)),N=v+1}let w=d.slice(N);return w.length>0&&S.push(O(w)),S.length===0?{type:"const",const:""}:S.length===1?S[0]:{type:"or",expr:S}}function b(d){function U(w,v){if(!u(w,v))throw new ue("TemplateLiteralParser: Index must point to open parens");let j=0;for(let F=v;F<w.length;F++)if(u(w,F)&&(j+=1),T(w,F)&&(j-=1),j===0)return[v,F];throw new ue("TemplateLiteralParser: Unclosed group parens in expression")}function N(w,v){for(let j=v;j<w.length;j++)if(u(w,j))return[v,j];return[v,w.length]}let S=[];for(let w=0;w<d.length;w++)if(u(d,w)){let[v,j]=U(d,w),F=d.slice(v,j+1);S.push(O(F)),w=j}else{let[v,j]=N(d,w),F=d.slice(v,j);F.length>0&&S.push(O(F)),w=j-1}return S.length===0?{type:"const",const:""}:S.length===1?S[0]:{type:"and",expr:S}}function O(d){return a(d)?O(s(d)):y(d)?R(d):I(d)?b(d):{type:"const",const:d}}c.Parse=O;function m(d){return O(d.slice(1,d.length-1))}c.ParseExact=m})(X||(i.TemplateLiteralParser=X={}));var Se=class extends k{};i.TemplateLiteralFiniteError=Se;var Y;(function(c){function t(s){throw new Se(s)}function u(s){return s.type==="or"&&s.expr.length===2&&s.expr[0].type==="const"&&s.expr[0].const==="0"&&s.expr[1].type==="const"&&s.expr[1].const==="[1-9][0-9]*"}function T(s){return s.type==="or"&&s.expr.length===2&&s.expr[0].type==="const"&&s.expr[0].const==="true"&&s.expr[1].type==="const"&&s.expr[1].const==="false"}function p(s){return s.type==="const"&&s.const===".*"}function a(s){return T(s)?!0:u(s)||p(s)?!1:s.type==="and"?s.expr.every(y=>a(y)):s.type==="or"?s.expr.every(y=>a(y)):s.type==="const"?!0:t("Unknown expression type")}c.Check=a})(Y||(i.TemplateLiteralFinite=Y={}));var me=class extends k{};i.TemplateLiteralGeneratorError=me;var Z;(function(c){function*t(s){if(s.length===1)return yield*s[0];for(let y of s[0])for(let I of t(s.slice(1)))yield`${y}${I}`}function*u(s){return yield*t(s.expr.map(y=>[...a(y)]))}function*T(s){for(let y of s.expr)yield*a(y)}function*p(s){return yield s.const}function*a(s){return s.type==="and"?yield*u(s):s.type==="or"?yield*T(s):s.type==="const"?yield*p(s):(()=>{throw new me("Unknown expression")})()}c.Generate=a})(Z||(i.TemplateLiteralGenerator=Z={}));var He;(function(c){function*t(a){let s=a.trim().replace(/"|'/g,"");return s==="boolean"?yield i.Type.Boolean():s==="number"?yield i.Type.Number():s==="bigint"?yield i.Type.BigInt():s==="string"?yield i.Type.String():yield(()=>{let y=s.split("|").map(I=>i.Type.Literal(I.trim()));return y.length===0?i.Type.Never():y.length===1?y[0]:i.Type.Union(y)})()}function*u(a){if(a[1]!=="{"){let s=i.Type.Literal("$"),y=T(a.slice(1));return yield*[s,...y]}for(let s=2;s<a.length;s++)if(a[s]==="}"){let y=t(a.slice(2,s)),I=T(a.slice(s+1));return yield*[...y,...I]}yield i.Type.Literal(a)}function*T(a){for(let s=0;s<a.length;s++)if(a[s]==="$"){let y=i.Type.Literal(a.slice(0,s)),I=u(a.slice(s));return yield*[y,...I]}yield i.Type.Literal(a)}function p(a){return[...T(a)]}c.Parse=p})(He||(i.TemplateLiteralDslParser=He={}));var ve=class{constructor(t){this.schema=t}Decode(t){return new Ae(this.schema,t)}};i.TransformDecodeBuilder=ve;var Ae=class{constructor(t,u){this.schema=t,this.decode=u}Encode(t){let u=P.Type(this.schema);return o.TTransform(u)?(()=>{let a={Encode:s=>u[i.Transform].Encode(t(s)),Decode:s=>this.decode(u[i.Transform].Decode(s))};return{...u,[i.Transform]:a}})():(()=>{let T={Decode:this.decode,Encode:t};return{...u,[i.Transform]:T}})()}};i.TransformEncodeBuilder=Ae;var In=0,we=class extends k{};i.TypeBuilderError=we;var xe=class{Create(t){return t}Throw(t){throw new we(t)}Discard(t,u){return u.reduce((T,p)=>{let{[p]:a,...s}=T;return s},t)}Strict(t){return JSON.parse(JSON.stringify(t))}};i.TypeBuilder=xe;var se=class extends xe{ReadonlyOptional(t){return this.Readonly(this.Optional(t))}Readonly(t){return{...P.Type(t),[i.Readonly]:"Readonly"}}Optional(t){return{...P.Type(t),[i.Optional]:"Optional"}}Any(t={}){return this.Create({...t,[i.Kind]:"Any"})}Array(t,u={}){return this.Create({...u,[i.Kind]:"Array",type:"array",items:P.Type(t)})}Boolean(t={}){return this.Create({...t,[i.Kind]:"Boolean",type:"boolean"})}Capitalize(t,u={}){return{...W.Map(P.Type(t),"Capitalize"),...u}}Composite(t,u){let T=i.Type.Intersect(t,{}),a=Re.ResolveKeys(T,{includePatterns:!1}).reduce((s,y)=>({...s,[y]:i.Type.Index(T,[y])}),{});return i.Type.Object(a,u)}Enum(t,u={}){if(f.IsUndefined(t))return this.Throw("Enum undefined or empty");let T=Object.getOwnPropertyNames(t).filter(s=>isNaN(s)).map(s=>t[s]),a=[...new Set(T)].map(s=>i.Type.Literal(s));return this.Union(a,{...u,[i.Hint]:"Enum"})}Extends(t,u,T,p,a={}){switch(z.Extends(t,u)){case l.Union:return this.Union([P.Type(T,a),P.Type(p,a)]);case l.True:return P.Type(T,a);case l.False:return P.Type(p,a)}}Exclude(t,u,T={}){return o.TTemplateLiteral(t)?this.Exclude(M.Resolve(t),u,T):o.TTemplateLiteral(u)?this.Exclude(t,M.Resolve(u),T):o.TUnion(t)?(()=>{let p=t.anyOf.filter(a=>z.Extends(a,u)===l.False);return p.length===1?P.Type(p[0],T):this.Union(p,T)})():z.Extends(t,u)!==l.False?this.Never(T):P.Type(t,T)}Extract(t,u,T={}){return o.TTemplateLiteral(t)?this.Extract(M.Resolve(t),u,T):o.TTemplateLiteral(u)?this.Extract(t,M.Resolve(u),T):o.TUnion(t)?(()=>{let p=t.anyOf.filter(a=>z.Extends(a,u)!==l.False);return p.length===1?P.Type(p[0],T):this.Union(p,T)})():z.Extends(t,u)!==l.False?P.Type(t,T):this.Never(T)}Index(t,u,T={}){return o.TArray(t)&&o.TNumber(u)?P.Type(t.items,T):o.TTuple(t)&&o.TNumber(u)?(()=>{let a=(f.IsUndefined(t.items)?[]:t.items).map(s=>P.Type(s));return this.Union(a,T)})():(()=>{let p=oe.Resolve(u),a=P.Type(t);return Ve.Resolve(a,p,T)})()}Integer(t={}){return this.Create({...t,[i.Kind]:"Integer",type:"integer"})}Intersect(t,u={}){if(t.length===0)return i.Type.Never();if(t.length===1)return P.Type(t[0],u);t.some(s=>o.TTransform(s))&&this.Throw("Cannot intersect transform types");let T=t.every(s=>o.TObject(s)),p=P.Rest(t),a=o.TSchema(u.unevaluatedProperties)?{unevaluatedProperties:P.Type(u.unevaluatedProperties)}:{};return u.unevaluatedProperties===!1||o.TSchema(u.unevaluatedProperties)||T?this.Create({...u,...a,[i.Kind]:"Intersect",type:"object",allOf:p}):this.Create({...u,...a,[i.Kind]:"Intersect",allOf:p})}KeyOf(t,u={}){return o.TRecord(t)?(()=>{let T=Object.getOwnPropertyNames(t.patternProperties)[0];return T===i.PatternNumberExact?this.Number(u):T===i.PatternStringExact?this.String(u):this.Throw("Unable to resolve key type from Record key pattern")})():o.TTuple(t)?(()=>{let p=(f.IsUndefined(t.items)?[]:t.items).map((a,s)=>i.Type.Literal(s.toString()));return this.Union(p,u)})():o.TArray(t)?this.Number(u):(()=>{let T=Re.ResolveKeys(t,{includePatterns:!1});if(T.length===0)return this.Never(u);let p=T.map(a=>this.Literal(a));return this.Union(p,u)})()}Literal(t,u={}){return this.Create({...u,[i.Kind]:"Literal",const:t,type:typeof t})}Lowercase(t,u={}){return{...W.Map(P.Type(t),"Lowercase"),...u}}Never(t={}){return this.Create({...t,[i.Kind]:"Never",not:{}})}Not(t,u){return this.Create({...u,[i.Kind]:"Not",not:P.Type(t)})}Null(t={}){return this.Create({...t,[i.Kind]:"Null",type:"null"})}Number(t={}){return this.Create({...t,[i.Kind]:"Number",type:"number"})}Object(t,u={}){let T=Object.getOwnPropertyNames(t),p=T.filter(I=>o.TOptional(t[I])),a=T.filter(I=>!p.includes(I)),s=o.TSchema(u.additionalProperties)?{additionalProperties:P.Type(u.additionalProperties)}:{},y=T.reduce((I,R)=>({...I,[R]:P.Type(t[R])}),{});return a.length>0?this.Create({...u,...s,[i.Kind]:"Object",type:"object",properties:y,required:a}):this.Create({...u,...s,[i.Kind]:"Object",type:"object",properties:y})}Omit(t,u,T={}){let p=oe.Resolve(u);return Q.Map(this.Discard(P.Type(t),["$id",i.Transform]),a=>{f.IsArray(a.required)&&(a.required=a.required.filter(s=>!p.includes(s)),a.required.length===0&&delete a.required);for(let s of Object.getOwnPropertyNames(a.properties))p.includes(s)&&delete a.properties[s];return this.Create(a)},T)}Partial(t,u={}){return Q.Map(this.Discard(P.Type(t),["$id",i.Transform]),T=>{let p=Object.getOwnPropertyNames(T.properties).reduce((a,s)=>({...a,[s]:this.Optional(T.properties[s])}),{});return this.Object(p,this.Discard(T,["required"]))},u)}Pick(t,u,T={}){let p=oe.Resolve(u);return Q.Map(this.Discard(P.Type(t),["$id",i.Transform]),a=>{f.IsArray(a.required)&&(a.required=a.required.filter(s=>p.includes(s)),a.required.length===0&&delete a.required);for(let s of Object.getOwnPropertyNames(a.properties))p.includes(s)||delete a.properties[s];return this.Create(a)},T)}Record(t,u,T={}){return o.TTemplateLiteral(t)?(()=>{let p=X.ParseExact(t.pattern);return Y.Check(p)?this.Object([...Z.Generate(p)].reduce((a,s)=>({...a,[s]:P.Type(u)}),{}),T):this.Create({...T,[i.Kind]:"Record",type:"object",patternProperties:{[t.pattern]:P.Type(u)}})})():o.TUnion(t)?(()=>{let p=qe.Resolve(t);if(o.TUnionLiteral(p)){let a=p.anyOf.reduce((s,y)=>({...s,[y.const]:P.Type(u)}),{});return this.Object(a,{...T,[i.Hint]:"Record"})}else this.Throw("Record key of type union contains non-literal types")})():o.TLiteral(t)?f.IsString(t.const)||f.IsNumber(t.const)?this.Object({[t.const]:P.Type(u)},T):this.Throw("Record key of type literal is not of type string or number"):o.TInteger(t)||o.TNumber(t)?this.Create({...T,[i.Kind]:"Record",type:"object",patternProperties:{[i.PatternNumberExact]:P.Type(u)}}):o.TString(t)?(()=>{let p=f.IsUndefined(t.pattern)?i.PatternStringExact:t.pattern;return this.Create({...T,[i.Kind]:"Record",type:"object",patternProperties:{[p]:P.Type(u)}})})():this.Never()}Recursive(t,u={}){f.IsUndefined(u.$id)&&(u.$id=`T${In++}`);let T=t({[i.Kind]:"This",$ref:`${u.$id}`});return T.$id=u.$id,this.Create({...u,[i.Hint]:"Recursive",...T})}Ref(t,u={}){return f.IsString(t)?this.Create({...u,[i.Kind]:"Ref",$ref:t}):(f.IsUndefined(t.$id)&&this.Throw("Reference target type must specify an $id"),this.Create({...u,[i.Kind]:"Ref",$ref:t.$id}))}Required(t,u={}){return Q.Map(this.Discard(P.Type(t),["$id",i.Transform]),T=>{let p=Object.getOwnPropertyNames(T.properties).reduce((a,s)=>({...a,[s]:this.Discard(T.properties[s],[i.Optional])}),{});return this.Object(p,T)},u)}Rest(t){return o.TTuple(t)&&!f.IsUndefined(t.items)?P.Rest(t.items):o.TIntersect(t)?P.Rest(t.allOf):o.TUnion(t)?P.Rest(t.anyOf):[]}String(t={}){return this.Create({...t,[i.Kind]:"String",type:"string"})}TemplateLiteral(t,u={}){let T=f.IsString(t)?Le.Create(He.Parse(t)):Le.Create(t);return this.Create({...u,[i.Kind]:"TemplateLiteral",type:"string",pattern:T})}Transform(t){return new ve(t)}Tuple(t,u={}){let[T,p,a]=[!1,t.length,t.length],s=P.Rest(t),y=t.length>0?{...u,[i.Kind]:"Tuple",type:"array",items:s,additionalItems:T,minItems:p,maxItems:a}:{...u,[i.Kind]:"Tuple",type:"array",minItems:p,maxItems:a};return this.Create(y)}Uncapitalize(t,u={}){return{...W.Map(P.Type(t),"Uncapitalize"),...u}}Union(t,u={}){return o.TTemplateLiteral(t)?M.Resolve(t):(()=>{let T=t;if(T.length===0)return this.Never(u);if(T.length===1)return this.Create(P.Type(T[0],u));let p=P.Rest(T);return this.Create({...u,[i.Kind]:"Union",anyOf:p})})()}Unknown(t={}){return this.Create({...t,[i.Kind]:"Unknown"})}Unsafe(t={}){return this.Create({...t,[i.Kind]:t[i.Kind]||"Unsafe"})}Uppercase(t,u={}){return{...W.Map(P.Type(t),"Uppercase"),...u}}};i.JsonTypeBuilder=se;var Ce=class extends se{AsyncIterator(t,u={}){return this.Create({...u,[i.Kind]:"AsyncIterator",type:"AsyncIterator",items:P.Type(t)})}Awaited(t,u={}){let T=p=>p.length>0?(()=>{let[a,...s]=p;return[this.Awaited(a),...T(s)]})():p;return o.TIntersect(t)?i.Type.Intersect(T(t.allOf)):o.TUnion(t)?i.Type.Union(T(t.anyOf)):o.TPromise(t)?this.Awaited(t.item):P.Type(t,u)}BigInt(t={}){return this.Create({...t,[i.Kind]:"BigInt",type:"bigint"})}ConstructorParameters(t,u={}){return this.Tuple([...t.parameters],{...u})}Constructor(t,u,T){let[p,a]=[P.Rest(t),P.Type(u)];return this.Create({...T,[i.Kind]:"Constructor",type:"Constructor",parameters:p,returns:a})}Date(t={}){return this.Create({...t,[i.Kind]:"Date",type:"Date"})}Function(t,u,T){let[p,a]=[P.Rest(t),P.Type(u)];return this.Create({...T,[i.Kind]:"Function",type:"Function",parameters:p,returns:a})}InstanceType(t,u={}){return P.Type(t.returns,u)}Iterator(t,u={}){return this.Create({...u,[i.Kind]:"Iterator",type:"Iterator",items:P.Type(t)})}Parameters(t,u={}){return this.Tuple(t.parameters,{...u})}Promise(t,u={}){return this.Create({...u,[i.Kind]:"Promise",type:"Promise",item:P.Type(t)})}RegExp(t,u={}){let T=f.IsString(t)?t:t.source;return this.Create({...u,[i.Kind]:"String",type:"string",pattern:T})}RegEx(t,u={}){return this.RegExp(t,u)}ReturnType(t,u={}){return P.Type(t.returns,u)}Symbol(t){return this.Create({...t,[i.Kind]:"Symbol",type:"symbol"})}Undefined(t={}){return this.Create({...t,[i.Kind]:"Undefined",type:"undefined"})}Uint8Array(t={}){return this.Create({...t,[i.Kind]:"Uint8Array",type:"Uint8Array"})}Void(t={}){return this.Create({...t,[i.Kind]:"Void",type:"void"})}};i.JavaScriptTypeBuilder=Ce;i.JsonType=new se;i.Type=new Ce});var De="messageLintRule.inlang.identicalPattern";var _e={en:"Identical pattern"},We={en:"Checks for identical patterns in different languages.  A message with identical wording in multiple languages can indicate that the translations are redundant or can be combined into a single message to reduce translation effort."};var ae=fn(Ye(),1),On=ae.Type.Object({ignore:ae.Type.Optional(ae.Type.Array(ae.Type.String({pattern:"[^*]",description:"All items in the array need quotaion marks at the end and beginning"}),{title:"DEPRECATED. Ignore paths",description:"Set a path that should be ignored."}))}),Ge={id:De,displayName:_e,description:We,settingsSchema:On,run:({message:c,report:t,settings:u})=>{let T=u[De],p=c.variants.find(s=>s.languageTag===u.sourceLanguageTag);if(p===void 0)return;let a=c.variants.filter(s=>s.languageTag!==u.sourceLanguageTag);for(let s of a){let y=Ze(p)===Ze(s),I=T?.ignore?.includes(bn(p.pattern));y&&!I&&t({messageId:c.id,languageTag:s.languageTag,body:{en:`Identical content found in language '${s.languageTag}' with message ID '${c.id}'.`}})}}},Ze=c=>{let t={...c,languageTag:void 0};return JSON.stringify(t)},bn=c=>c.filter(t=>t.type==="Text").map(t=>t.value).join("");var mn=Ge;export{mn as default};
