import { useEffect } from "react";
import { useServe } from "../utils/customHook.js";

// eslint-disable-next-line react/prop-types
  function Images({ url, params = {}, captions }) {
    const { images, e, getImages } = useServe();

    useEffect(() => {
      getImages(url, params);    
      // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [url,params]);

  return (
    <>
      {e ? (
        <img src={"qweasd"} alt={captions || "not found images"} />
      ) : (
        <img
          src={images}
          alt={captions || "not found images"}
          className="w-full h-full  object-fill"
        />
      )}
    </>
  );
}

export default Images;
