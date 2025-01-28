import http from "k6/http";
import { group, check } from "k6";
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.6.0/index.js';

const BASE_URL = `${__ENV.PLANS_ENDPOINT}`;

const getURL = () => BASE_URL;

export const listPlans = () => {
  const req = http.get(getURL(), {headers: {"Content-Type": "application/json", "Accept": "application/json"}});
  return req;
}

export function testListPlans() {
  group("list plans", () => {
    let req = listPlans();
    check(req, { "Successful listing.": (r) => r.status === 200 });
  });
}

export const generatePlan = () => {
  const randomU = Math.floor(Math.random() * 101);

  const letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
  const randomLetters = letters[Math.floor(Math.random() * 26)] + letters[Math.floor(Math.random() * 26)];

  const obj = {
    id: uuidv4(),
    name: randomLetters + randomU,
    price: randomU * randomU ,
    description: "Plan"
  };

  return obj
};

function addPlan() {
  const req = http.post(getURL(), JSON.stringify(generatePlan()), {headers: {"Content-Type": "application/json", "Accept": "application/json"}});
  return req;
}

export function testAddPlan() {
  group("add plan", () => {
    let req = addPlan();
    check(req, { "Plan created.": (r) => r.status === 200 });// bug: service should return 201
  });
}

export function testAddPlans(qty){
  group("add " + qty + " plans", () => {
    for (let index = 0; index < qty; index++) {
      let req = addPlan();
      check(req, { "Plan created.": (r) => r.status === 200 }); // bug: service return should 201
    }
  });
}