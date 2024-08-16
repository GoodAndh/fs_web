/* eslint-disable no-prototype-builtins */
import { useState } from "react";
import Form from "../components/Form.jsx";
import Navbar from "../components/Navbar.jsx";
import { usePostJson } from "../utils/customHook.js";
import { Link } from "react-router-dom";

function Register() {
  const initialForm = {
    username: "",
    email: "",
    password: "",
    cpassword: "",
  };
  const [form, setForm] = useState(initialForm);
  const { data, error, call } = usePostJson("signup");

  const onChange = (id, value) => {
    setForm({ ...form, [id]: value });
  };

  const onSubmit = async (e) => {
    e.preventDefault();
    await call({ form: form }); 
  };

  return (
    <>
      <Navbar></Navbar>
      <div className="mt-60 max-w-lg m-3 mx-auto p-5 border-2 border-slate-600 rounded-xl shadow-md">
        {data && (
          <div
            className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded-md relative"
            role="alert"
          >
            <strong className="font-bold">Akun berhasil dibuat!</strong>
            <p className="block ml-4 sm:inline">
              Kembali ke halaman{" "}
              <Link to="/signin" className="font-bold">
                Login
              </Link>
            </p>
          </div>
        )}
        <Form
          id="username"
          type="text"
          title="Username"
          inputChange={onChange}
          validate={ error?.data?.Username || ""}
        ></Form>
        <Form
          id="email"
          type="text"
          title="Email"
          inputChange={onChange}
          validate={error?.data?.Email || ""}
        ></Form>
        <Form
          id="password"
          type="password"
          title="Password"
          inputChange={onChange}
          validate={ error?.data?.Password || ""}
        ></Form>
        <Form
          id="vpassword"
          type="password"
          title="Validate Password"
          inputChange={onChange}
          validate={error?.data?.VPassword || ""}
        ></Form>
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

export default Register;
