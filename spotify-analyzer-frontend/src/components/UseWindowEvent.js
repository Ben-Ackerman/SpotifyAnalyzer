import { useEffect } from 'react';

// This is a custom hook that takes in an an event to add a listener to, with a callback of callback
export const useWindowEvent = (event, callback) => {
    useEffect(() => {
        window.addEventListener(event, callback);
        return () => {window.removeEventListener(event, callback)};
    }, [event, callback]);
};