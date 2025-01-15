import http from "k6/http";
import { group, check } from "k6";
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.6.0/index.js';
import { listUsers } from "./svc-users.js";
import { listPlans } from "./svc-plans.js";

const BASE_URL = `${__ENV.SUBSCRIPTIONS_ENDPOINT}`;

const getURL = () => BASE_URL;

export const listSubscriptions = () => {
  const req = http.get(getURL(), {headers: {"Content-Type": "application/json", "Accept": "application/json"}});
  return req;
}

export function testListSubscriptions() {
  group("list subscriptions", () => {
    let req = listSubscriptions();
    check(req, { "Successful listing.": (r) => r.status === 200 });
  });
}

const generateSubscription = (userId, planId) => {
  const obj = {
    id: uuidv4(),
    user_id: userId,
    plan_id: planId
  };
  return obj
};

const addSubscription = (userId, planId) => {
  const req = http.post(getURL(), JSON.stringify(generateSubscription(userId, planId)), {headers: {"Content-Type": "application/json", "Accept": "application/json"}});
  return req;
}

export function testAddSubscription(userId, planId) {
  group("add Subscription", (userId, planId) => {
    let req = addSubscription(userId, planId);
    check(req, { "Successful listing.": (r) => r.status === 201 });
  });
}

export function testAddSubscriptions(qty){
  group("add " + qty + " subscriptions", () => {
    for (let i = 0; i < qty; i++) {
        
      const usersResponse = listUsers()
      if (usersResponse.status !== 200) {
        console.error("Failed to list users:", usersResponse.error_code);
        continue;
      }

      try {
        const users = usersResponse.json();
        const randomUser = users[Math.floor(Math.random() * users.length)];
        const userId = randomUser.id;

        const plansResponse = listPlans()
        if (plansResponse.status !== 200) {
          console.error("Failed to list plans:", plansResponse.error_code);
          continue;
        }
        const plans = plansResponse.json();
        const randomPlan = plans[Math.floor(Math.random() * plans.length)];
        const planId = randomPlan.id;
      
        const response = addSubscription(userId, planId);
        check(response, {
          "Subscription created.": (r) => r.status === 200, // bug: service return should 201
        });

        if (response.status !== 200) {
          console.error("Failed to add a subscription:", response.error_code);
        }

      } catch (error) {
        console.error("Error processing subscriptions:", error);
      }
    }
  });
}

