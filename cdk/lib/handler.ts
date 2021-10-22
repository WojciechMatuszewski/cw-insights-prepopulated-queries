export const handler = async () => {
  console.log("Start");
  return {
    statusCode: 200,
    body: JSON.stringify("Hi")
  };
};
