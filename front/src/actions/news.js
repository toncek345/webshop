import { ServerIp } from '../constants';

export const getNews = () => (dispatch) => {
  dispatch({ type: 'GET_NEWS_PENDING' });
  fetch(`${ServerIp}/news`, {
    method: 'GET',
  }).then((resp) => {
    if (resp.status < 400) {
      resp.json().then((data) => {
        dispatch({
          type: 'GET_NEWS_DONE',
          payload: data,
        });
      });
    } else {
      dispatch({
        type: 'GET_NEWS_FAIL',
        payload: 'Error loading news.',
      });
    }
  });
};
