/* eslint-disable react/prop-types */
function Form({ inputChange, id, validate, type, title }) {
  const baseClass = `mt-2 px-3 py-2 shadow-md rounded-xl w-full block text-md cursor-pointer `;
  return (
    <>
      <BaseForm
        inputChange={inputChange}
        id={id}
        validate={validate}
        type={type}
        title={title}
        customClassesForm={baseClass}
      />
    </>
  );
}

export default Form;

export function BaseForm({
  inputChange,
  id,
  validate,
  type,
  title,
  customClassesForm,
  newLabel,
}) {
  const InvalidClass = "outline outline-red-800";
  const invalidMsgClass = "text-sm m-3 font-semibold text-pink-800 ";

  const onInputChange = (e) => {
    const { id, value, files, type } = e.target;
    if (type === "number") {
      const numericValue = parseFloat(value);

      if (isNaN(numericValue) || numericValue < 0) {
        return;
      }
    }
    files ? inputChange(id, value, files) : inputChange(id, value);
  };
  return (
    <>
      <div>
        <span className=" m-3 font-semibold block hover:cursor-pointer ">
          {title}
        </span>
        {type === "textarea" ? (
          <textarea
            onChange={onInputChange}
            id={id}
            className={`${customClassesForm} ${
              validate !== "" ? InvalidClass : ""
            }`}
          />
        ) : (
          <input
            onChange={onInputChange}
            id={id}
            type={type}
            {...(type === "number" && { min: 0 })}
            className={` ${customClassesForm} ${
              validate !== "" ? InvalidClass : ""
            }`}
            {...(type === "file" && { accept: "image/*" })}
          />
        )}

        {validate !== "" && <p className={invalidMsgClass}>{validate}</p>}
      </div>
      {newLabel && newLabel}
    </>
  );
}
