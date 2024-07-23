import { useState, useEffect } from "react";
import { GoChevronLeft, GoChevronRight } from "react-icons/go";
import { useGetJson } from "../utils/customHook.js";
import Images from "./Images.jsx";

// eslint-disable-next-line react/prop-types
function ImageSlider({ itemID }) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [urlImage, setUrlImage] = useState([]);
  const { error, data } = useGetJson(`product/image/${itemID}`);

  useEffect(() => {
    if (data) {
      setUrlImage(data.data);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data]);

  const prevSlide = () => {
    const isFirstSlide = currentIndex === 0;
    const newIndex = isFirstSlide ? urlImage.length - 1 : currentIndex - 1;
    setCurrentIndex(newIndex);
  };
  const nextSlide = () => {
    const isLastSlide = currentIndex === urlImage.length - 1;
    const newIndex = isLastSlide ? 0 : currentIndex + 1;
    setCurrentIndex(newIndex);
  };

  return (
    <>
      {error ? (
        <div>ada error</div>
      ) : (
        <div className=" bg-center overflow-hidden  relative group">
          <div className="w-full h-full rounded-2xl bg-center duration-500">
            {urlImage[currentIndex]?.url ? (
              <Images
                url="product/serveimg"
                params={{ url: `${urlImage[currentIndex]?.url}` }}
                captions={`${urlImage[currentIndex]?.captions}`}
              />
            ) : (
              "is loading..."
            )}
          </div>
          {/* left arrow */}
          <div className="hidden group-hover:block absolute top-[50%] -translate-x-0 translate-y-[-50%] left-5 text-2xl cursor-pointer">
            <GoChevronLeft onClick={prevSlide} size={50} />
          </div>
          {/* right arrow */}
          <div className="hidden group-hover:block absolute top-[50%] -translate-x-0 translate-y-[-50%] right-5 text-2xl cursor-pointer">
            <GoChevronRight onClick={nextSlide} size={50} />
          </div>
        </div>
      )}
    </>
  );
}

export default ImageSlider;
