import { ServerIp } from '../constants';

export const getProducts = () => (dispatch) => {
  dispatch({ type: 'GET_PRODUCTS_PENDING' });
  fetch(`${ServerIp}/product`, {
    method: 'GET',
  }).then((resp) => {
    if (resp.status < 400) {
      resp.json().then((data) => {
        dispatch({
          type: 'GET_PRODUCTS_DONE',
          payload: data,
        });
      });
    } else {
      dispatch({
        type: 'GET_PRODUCTS_FAIL',
        payload: 'Error loading products.',
      });
    }
  });
};

export function increaseNumberOfItem(item, times) {
  return {
    type: 'INCREASE_NUMBER_OF_ITEM',
    payload: {
      item,
      times,
    },
  };
}

export function addToCart(item) {
  return {
    type: 'ADD_TO_CART',
    payload: item,
  };
}

export function removeFromCart(item) {
  return {
    type: 'REMOVE_FROM_CART',
    payload: item,
  };
}
