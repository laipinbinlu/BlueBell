import { message } from 'antd';
import { ApiResponse } from '../types';

export const handleApiError = (error: any, fallbackMessage: string) => {
  console.error(fallbackMessage, error);
  message.error(error.response?.data?.msg || fallbackMessage);
};

export const handleApiResponse = <T>(
  response: ApiResponse<T>,
  onSuccess: (data: T) => void,
  fallbackErrorMessage: string
) => {
  if (response.code === 1000) {
    onSuccess(response.data);
    return true;
  } else {
    message.error(response.msg || fallbackErrorMessage);
    return false;
  }
}; 