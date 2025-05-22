export const get = async <T>(
  path: string,
  options?: RequestInit,
): Promise<ApiResult<T>> => {
  const baseUrl = import.meta.env.VITE_REACT_APP_API_BASE_URL;
  const url: string = `${baseUrl}/api/${path}`;
  try {
    const response: Response = await fetch(url, options ? options : {});
    return handleResponse(response);
  } catch (e: any) {
    let message = "Unknown error occurred";
    if (e instanceof Error) {
      message = e.message;
    }
    return {
      ok: false,
      data: {
        message: message,
      },
    };
  }
};

export const post = async <P, T>(
  path: string,
  data: P,
  options?: RequestInit,
): Promise<ApiResult<T>> => {
  const baseUrl = import.meta.env.VITE_REACT_APP_API_BASE_URL;
  const url: string = `${baseUrl}/api/${path}`;
  try {
    const response: Response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
      ...options,
    });
    return handleResponse(response);
  } catch (e: any) {
    let message = "Unknown error occurred";
    if (e instanceof Error) {
      message = e.message;
    }
    return {
      ok: false,
      data: {
        message: message,
      },
    };
  }
};

export const del = async <T>(
  path: string,
  options?: RequestInit,
): Promise<ApiResult<T>> => {
  const baseUrl = import.meta.env.VITE_REACT_APP_API_BASE_URL;
  const url: string = `${baseUrl}/api/${path}`;
  try {
    const response: Response = await fetch(url, {
      method: "DELETE",
      ...options,
    });
    return handleResponse(response);
  } catch (e: any) {
    let message = "Unknown error occurred";
    if (e instanceof Error) {
      message = e.message;
    }
    return {
      ok: false,
      data: {
        message: message,
      },
    };
  }
}

const handleResponse = async <T>(response: Response): Promise<ApiResult<T>> => {
  const isJson = response.headers.get("content-type")?.includes(
    "application/json",
  );
  const data = isJson ? await response.json() : null;
  if (response.ok) {
    return { ok: true, data };
  }
  const error = (data && data.message) || response.statusText;
  return {
    ok: false,
    data: {
      message: error,
      code: data?.code,
    },
  };
};

export type ApiResult<T> = ApiSuccessResult<T> | ApiErrorResult;

export type ApiSuccessResult<T> = {
  data: T;
  ok: true;
};

export type ApiErrorResult = {
  data: { code?: string; message: string };
  ok: false;
};
