import * as fs from 'fs'
import * as  path from 'path';


var readFileTime
var calculateInputTime
function getElapsedTimeInSeconds(start) {
  const diff = process.hrtime(start);
  // Convert to seconds
  return Number((diff[0] + diff[1] / 1e9).toFixed(2)); // seconds + fractional part of seconds
}

function readData() {
  const startTime = process.hrtime();
  const jsonData0 = fs.readFileSync(path.resolve('data0.json'));
  const jsonData1 = fs.readFileSync(path.resolve('data1.json'));

  // Parse it into an object
  const data0 = JSON.parse(jsonData0);
  const data1 = JSON.parse(jsonData1);
  const allData = { pairs: [...data0.pairs, ...data1.pairs] }

  readFileTime = getElapsedTimeInSeconds(startTime);

  return allData
}
function degreesToRadians(degrees) {
  return degrees * (Math.PI / 180);
}

function haversineOfDegress(x0, y0, x1, y1, r) {

  const dY = degreesToRadians(y1 - y0)
  const dX = degreesToRadians(x1 - x0)
  y0 = degreesToRadians(y0)
  y1 = degreesToRadians(y1)

  const rootTerm = Math.pow(Math.sin(dY / 2),2) + Math.cos(y0) * Math.cos(y1) * Math.pow(Math.sin(dX / 2),2)
  const result = 2 * r * Math.asin(Math.sqrt(rootTerm))

  return result
}

function calculateHaversine(data) {
  const startTime = process.hrtime();
  const earthRadiuskm = 6371
  let sum = 0
  let count = 0
  for (let item of data.pairs) {
    sum += haversineOfDegress(item.x0, item.y0, item.x1, item.y1, earthRadiuskm)
    count += 1
  }
  const avg = sum / count

  calculateInputTime = getElapsedTimeInSeconds(startTime);

  return  {avg, count}
}

function haversine() {
  const data = readData()
  const {avg, count} = calculateHaversine(data)

  console.log("Result: ", avg)
  console.log("Input = ", readFileTime, " seconds")
  console.log("Math = ", calculateInputTime, " seconds")
  console.log("Total = ", readFileTime + calculateInputTime, " seconds")
  console.log("Throughput = ", count / (readFileTime+ calculateInputTime), " haversines/second")
}


haversine()
