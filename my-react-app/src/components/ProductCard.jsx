import ImageSlider from "./ProductImageSlider.jsx";
import { useGetJson } from "../utils/customHook.js";
import { Link } from "react-router-dom";


function ProductCard() {
  const { error, data } = useGetJson("product");

  return (
    <>
      {error ? (
        <div>ada erorr</div>
      ) : (
        <>
          {data &&
            data.data &&
            data.data.map((item) => {
              return (
                <div
                  key={item.id}
                  className="flex items-center justify-center  container mx-auto mt-12"
                >
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    <div className="rounded-xl shadow-lg  bg-slate-300 ">
                      <div className="m-2 p-1">
                        <ImageSlider itemID={item.id} />
                        <br />
                        <h5 className="text-2xl md:text-3xl m-1">
                          {item.name}
                        </h5>
                        <p className="text-md font-medium m-1">
                          Rp.{item.price}
                        </p>
                        <div className="my-4">
                          <Link
                            to={`/p/${item.name}`}
                            className="my-5 w-full p-2 bg-transparent hover:bg-green-500 text-slate-700 font-semibold hover:text-white py-2 px-4 border border-slate-500 hover:border-transparent rounded-xl"
                          >
                            Selengkapnya
                          </Link>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}
        </>
      )}
    </>
  );
}

export default ProductCard;

// <div className="flex items-center justify-center min-h-screen container mx-auto mt-12">
//             <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
//               <div className="rounded-xl shadow-lg  bg-slate-300 ">
//                 <div className="m-2 p-1">
//                   <div className="rounded-xl overflow-hidden ">
//                     <h5>
//                       Lorem ipsum, dolor sit amet consectetur adipisicing elit.
//                       Nam corrupti neque eius tenetur earum repellendus eum ab
//                       aliquam consectetur quis illum cum quisquam nesciunt
//                       molestias sint, alias nihil ut praesentium quae ducimus
//                       minus nostrum quam modi? Illum quaerat nesciunt
//                       saepe?IMAGE
//                     </h5>
//                     <img src={"images"} alt="ERROR BRO SORRY MY BAD" />
//                   </div>
//                   <br />
//                   <h5 className="text-2xl md:text-3xl m-1">
//                     Lorem ipsum dolor sit amet.
//                   </h5>
//                   <p className="text-md font-medium m-1">
//                     Lorem ipsum dolor sit amet consectetur.
//                   </p>
//                   <div className="my-4">
//                     <a
//                       href="#"
//                       className="my-5 w-full p-2 bg-transparent hover:bg-green-500 text-slate-700 font-semibold hover:text-white py-2 px-4 border border-slate-500 hover:border-transparent rounded-xl"
//                     >
//                       Selengkapnya
//                     </a>
//                   </div>
//                 </div>
//               </div>
//             </div>
//           </div>
