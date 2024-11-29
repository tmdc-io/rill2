export function fetchTokenFromLocalStorage() {
  let token;
  const keys = Object.keys(localStorage).filter(key => key.startsWith("modern-oidc"));
  if(keys?.length>0){
    const storeValue = localStorage.getItem(keys[0]) ? JSON.parse(localStorage.getItem(keys[0])) : {};
    token = storeValue?.access_token;
  }
  return token;
}