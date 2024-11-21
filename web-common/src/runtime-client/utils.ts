export function fetchTokenFromLocalStorage() {
  let token;
  const keys = Object.keys(localStorage).filter(key => key.startsWith("modern-oidc"));
  console.warn("keys========>",keys);
  if(keys?.length>0){
    const storeValue = localStorage.getItem(keys[0]) ? JSON.parse(localStorage.getItem(keys[0])) : {};
    console.warn("keys========>",storeValue);
    token = storeValue?.access_token;
  }
  console.warn("keys========>",token);
  return token;
}