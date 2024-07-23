// eslint-disable-next-line react/prop-types
function Form({ inputChange, id, validate, type, title }) {
  const InvalidClass = "border-2 border-red-800 outline-red-800 peer";
  const invalidMsgClass =
    "text-sm m-3 font-semibold invisible text-pink-800 peer-valid:visible";

  const onInputChange = (e) => {
    const { id, value, files } = e.target;
    files ? inputChange(id, value, files) : inputChange(id, value);
    // files&&console.log("files:",files[0].name," values:",value);
  };
  return (
    <>
      <label htmlFor={id}>
        <span className="m-3 font-semibold block hover:cursor-pointer after:content-['*'] after:text-pink-500">
          {title}
        </span>
        <input
          onChange={onInputChange}
          id={id}
          type={type}
          className={`mt-2 px-3 py-2 shadow-md rounded-xl w-full block text-md hover:cursor-pointer  ${
            validate !== "" ? InvalidClass : ""
          }`}
          {...(type === "file" && { accept: "image/*" })}
        />
        {validate !== "" && <span className={invalidMsgClass}>{validate}</span>}
      </label>
    </>
  );
}

export default Form;
