export function fetchTokenFromLocalStorage() {
    const storeValue = localStorage.getItem("userInfo") ? JSON.parse(localStorage.getItem("userInfo")) : {};
    const TOKEN = storeValue?.accessToken;
    console.log('$$$$$$$$$$$$$$$ fetchWrapper TOKEN', TOKEN);
    return TOKEN;
  }