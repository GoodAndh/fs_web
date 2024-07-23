// import React from "react";
import Login from "./pages/login.jsx";
import Home from "./pages/home.jsx";
import Register from "./pages/register.jsx";
import Profile from "./pages/profile.jsx";
import ProductPage from "./pages/productPage.jsx";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

// eslint-disable-next-line react-refresh/only-export-components
function App() {
  const router = createBrowserRouter([
    {
      path: "*",
      element: <div>Error Wildcard</div>,
    },
    {
      path: "/signin",
      element: <Login />,
    },
    {
      path: "/",
      element: <Home />,
    },
    {
      path: "/signup",
      element: <Register />,
    },
    {
      path: "/profile",
      element: <Profile />,
    },
    {
      path: "/p/:nm",
      element: <ProductPage />,
    },
  ]);

  return (
    <>
      <RouterProvider router={router} />
    </>
  );
}

export default App;
