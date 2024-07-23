import axios from "axios";

const instance = axios.create({
  baseURL: "http://localhost:3000/api/",
  withCredentials: true,
});

instance.interceptors.response.use(
  function (response) {
    return response;
  },
  function (error) {
    return Promise.reject(error);
  }
);

instance.interceptors.request.use(
  function (config) {
    return config;
  },
  function (error) {
    return Promise.reject(error);
  }
);

export default instance;

export async function getJson(url, params, config = {}) {
  try {
    const response = await instance.get(url, {
      ...config,
      params: params,
    });
    return response.data;
  } catch (error) {
    return error.message;
  }
}

export async function postJson(url, data, headers) {
  try {
    const response = await instance.post(url, data, { headers: headers });
    if (response.status!==200) {
      throw new Error(response.response)
    }
    return response.data;
  } catch (error) {
    return error.response;
  }
}
