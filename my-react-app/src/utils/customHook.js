import { useEffect, useState } from "react";
import instance from "./instance.js";

export function useServe() {
  const [images, setImages] = useState(null);
  const [e, setError] = useState(null);

  const getImages = async (url, params) => {
    try {
      const response = await instance.get(url, {
        responseType: "blob",
        params: params,
      });
      const ObjUrl = URL.createObjectURL(response.data);
      setImages(ObjUrl);
    } catch (error) {
      setError(error.response);
    }
  };
  return { images, e, getImages };
}

export function useGetJson(url, params, config = {}) {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const call = async function () {
    setLoading(true);
    try {
      const response = await instance.get(url, {
        ...config,
        params: params,
      });
      if (response.status === 200) {
        setData(response.data);
      }
    } catch (e) {
      setError(e.response.data);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    call();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [url, params]);

  return { data, error, loading };
}

export function usePostJson(url) {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const call = async function ({ form, headers, params, fallbackFunc }) {
    setLoading(true);
    try {
      const queryURL = params
        ? `${url}?${new URLSearchParams(params).toString()}`
        : url;
      const response = await instance.post(queryURL, form, {
        headers: headers,
      });
      setData(response.data);
    } catch (e) {
      setError(e.response.data);
      if (fallbackFunc) {
        await fallbackFunc({ form, params });
      }
    } finally {
      setLoading(false);
    }
  };

  const resetResponse = () => {
    setError(null);
  };

  return { data, error, loading, call, resetResponse };
}
