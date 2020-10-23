const defaultState = {
  allProducts: null,
  loading: null,
  error: null,

  cartItems: [],
};

export default function Products(state = defaultState, action) {
  switch (action.type) {
    case "GET_PRODUCTS_PENDING":
      return {
        ...state,
        loading: true,
      };
    case "GET_PRODUCTS_DONE":
      return {
        ...state,
        loading: false,
        allProducts: action.payload,
      };
    case "GET_PRODUCTS_FAIL":
      return {
        ...state,
        loading: false,
        error: action.payload,
      };
    case "ADD_TO_CART":
      return {
        ...state,
        cartItems: state.cartItems.concat({
          ...action.payload,
          numberOfItems: 1,
        }),
      };
    case "INCREASE_NUMBER_OF_ITEM":
      return {
        ...state,
        cartItems: state.cartItems.map((item) => {
          if (item.Id === action.payload.item.Id) {
            return {
              ...item,
              numberOfItems: action.payload.times,
            };
          }
          return { ...item };
        }),
      };
    case "REMOVE_FROM_CART":
      return {
        ...state,
        cartItems: state.cartItems.filter(
          (item) => item.Id !== action.payload.Id
        ),
      };
    default:
      return state;
  }
}
