import { useState, useEffect } from "react";
import Navbar from "../components/Navbar.jsx";
import Form from "../components/Form.jsx";
import { useNavigate } from "react-router-dom";
import { usePostJson } from "../utils/customHook.js";

function Login() {
  const [form, setForm] = useState({
    username: "",
    password: "",
  });
  const navigate = useNavigate();
  const { data, error, call } = usePostJson("signin");
  const onChange = (id, value) => {
    setForm({ ...form, [id]: value });
  };
  const onSubmit = async (e) => {
    e.preventDefault();
    await call({ form: form });
  };

  useEffect(() => {
    if (data) {
      navigate("/");
      window.location.reload();
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data]);

  return (
    <>
      <Navbar></Navbar>
      <div className="mt-80 max-w-lg m-3 mx-auto p-5 border-2 border-slate-600 rounded-xl shadow-md">
        <Form
          id="username"
          type="text"
          title="Username"
          inputChange={onChange}
          validate={error?.message || ""}
        ></Form>
        <Form
          id="password"
          type="password"
          title="Password"
          inputChange={onChange}
          validate={error?.message || ""}
        ></Form>
        <p className="mt-4 text-center font-semibold">
          belum punya akun?
          <a href="/signup" className="ml-1 text-green-500 hover:underline">
            Buat Akun
          </a>
        </p>
        <button
          onClick={onSubmit}
          className="my-5 w-full p-2 bg-transparent hover:bg-green-500 text-slate-700 font-semibold hover:text-white py-2 px-4 border border-slate-500 hover:border-transparent rounded-xl"
        >
          Submit
        </button>{" "}
      </div>
    </>
  );
}

export default Login;
