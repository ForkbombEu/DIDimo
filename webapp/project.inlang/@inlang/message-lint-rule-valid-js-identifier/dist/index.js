var a="messageLintRule.inlang.validJsIdentifier";var t={en:"Valid JS Identifier"},s={en:"Make sure that all message IDs are valid JavaScript identifiers."};var r=["break","case","catch","class","const","continue","debugger","default","delete","do","else","export","extends","false","finally","for","function","if","import","in","instanceof","new","null","return","super","switch","this","throw","true","try","typeof","var","void","while","with","let","static","yield","await","enum","implements","interface","package","private","protected","public"];function l(e){return!r.includes(e)&&o(e)}function o(e){if(e.trim()!==e)return!1;try{new Function(e,"var "+e)}catch{return!1}return!0}var c={id:a,displayName:t,description:s,run:({message:e,settings:n,report:d})=>{l(e.id)||d({messageId:e.id,languageTag:n.sourceLanguageTag,body:{en:`The message ID '${e.id}' is not a valid javascript identifier.`}})}},v=c;export{v as default};
