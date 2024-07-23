
// eslint-disable-next-line react/prop-types
function Modal({ isOpen, title, content, onCloseButton }) {
  if (!isOpen) {
    return null;
  }
  return (
    <>
      <div
        className={`fixed inset-0 flex items-center justify-center z-50 transition-opacity ${
          isOpen ? "opacity-1000" : "opacity-0"
        }`}
      >
        <div className="fixed inset-0 bg-gray-500 opacity-50"></div>
        <div className="bg-white p-6 rounded-lg shadow-lg z-10 max-w-sm mx-auto transition-transform transform scale-100">
          <span className="font-semibold block  after:content-['*'] after:text-pink-500">
            {title || "Insert Title"}
          </span>{" "}
          {content || "empty content"}
          <div className="mt-[25px]">
           {onCloseButton}
          </div>
        </div>
      </div>
    </>
  );
}

export default Modal;
