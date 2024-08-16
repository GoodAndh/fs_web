import Modal from "./Modal.jsx";
import { Link } from "react-router-dom";

// eslint-disable-next-line react/prop-types
function ModalLoginErrorMsg({ closeButton }) {
  return (
    <>
      <Modal
        title="Sesi anda telah habis silahkan login kembali"
        isOpen={true}
        content={
          <div className="">
            <p className="mt-[25px] text-red-700">
              Kembali ke halaman{" "}
              <Link to="/signin" className="font-bold">
                Login
              </Link>
              <span className="">
                {" "}
                atau klik{" "}
                <Link to="/signin" className="font-bold">
                  disini
                </Link>
              </span>
            </p>
          </div>
        }
        onCloseButton={closeButton}
      />
    </>
  );
}

export default ModalLoginErrorMsg;
