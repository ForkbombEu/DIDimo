var i="messageLintRule.inlang.missingTranslation";var t={en:"Missing translation"},g={en:"Checks for missing variants for a specific languageTag.  If a variant exists for the sourceLanguageTag but is missing for a listed languageTag, it is likely that the message has not been translated for this languageTag yet."};var l={id:i,displayName:t,description:g,run:({message:e,settings:n,report:r})=>{let o=n.languageTags.filter(a=>a!==n.sourceLanguageTag);for(let a of o)(e.variants.filter(m=>m.languageTag===a)??[]).length||r({messageId:e.id,languageTag:a,body:{en:`Message with id '${e.id}' has a missing variant for language tag '${a}'.`}})}};var h=l;export{h as default};
