import http from "k6/http";
import { group, check } from "k6";
// import { faker } from "https://esm.sh/@faker-js/faker@v8.4.1"
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.6.0/index.js';

const BASE_URL = `${__ENV.USERS_ENDPOINT}`;

const getURL = () => BASE_URL;

export const listUsers = () => {
  const req = http.get(getURL(), {headers: {"Content-Type": "application/json", "Accept": "application/json"}});
  return req;
}

const generateUser = () => {
  const obj = {
    id: uuidv4(),
    name: "teste" + uuidv4(), //faker.person.firstName(),
    email: "teste" + uuidv4() + "@test.com" //faker.internet.email()
  };
  return obj
};

const addUser = () => {
  const req = http.post(getURL(), JSON.stringify(generateUser()), {headers: {"Content-Type": "application/json", "Accept": "application/json"}});
  return req;
}

export function testListUsers() {
  group("list users", () => {
    const req = listUsers();
    check(req, { "Successful listing.": (r) => r.status === 200 });
  });
}

export function testAddUser() {
  group("add user", () => {
    let req = addUser();
    check(req, { "User created.": (r) => r.status === 200 }); // bug: service return should 201
  });
}

export function testAddUsers(qty){
  group("add " + qty + " users", () => {
    for (let index = 0; index < qty; index++) {
      let req = addUser();
      check(req, { "User created.": (r) => r.status === 200 }); // bug: service return should 201
    }
  });
};

