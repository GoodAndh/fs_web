import { createContext, useState, useEffect, useContext } from "react";
import { useGetJson } from "./customHook.js";

const AuthContext = createContext();

// eslint-disable-next-line react/prop-types
function AuthProvider({ children }) {
  const [isAuth, setIsAuth] = useState(false);
  const { error, data } = useGetJson("user/validate");

  useEffect(() => {
    if (data) {
      setIsAuth(true);
    }
    if (error) {
      setIsAuth(false);
    }
  }, [data, error]);

  return (
    <AuthContext.Provider value={{ isAuth }}>{children}</AuthContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export const useAuth = () => useContext(AuthContext);

export default AuthProvider;
