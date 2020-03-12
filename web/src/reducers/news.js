const defaultState = {
  allNews: null,
  loading: null,
  error: null,
};

export default function News(state = defaultState, action) {
  switch (action.type) {
    case 'GET_NEWS_PENDING':
      return {
        ...state,
        loading: true,
      };
    case 'GET_NEWS_DONE':
      return {
        ...state,
        loading: false,
        allNews: action.payload,
      };
    case 'GET_NEWS_FAIL':
      return {
        ...state,
        loading: false,
        error: action.payload,
      };
    default:
      return state;
  }
}
