from selenium import webdriver
import random
import string


def register_user_test(test_email):
    chrome = webdriver.Chrome()
    chrome.get('http://localhost:8001/')

    # Wait for any redirects
    chrome.implicitly_wait(5) # seconds

    # Login with Multi-Factor
    chrome.find_element_by_id("login-button").click()

    # Select: Register new user
    chrome.find_elements_by_xpath("//*[text()='Register new user']")[0].click()

    # login input id = login
    username = chrome.find_element_by_id("input-12")
    username.clear()
    username.send_keys(test_email)

    # password input id = password
    password = chrome.find_element_by_id("password")
    password.clear()
    password.send_keys("secret-password")

    # click register button
    chrome.find_element_by_class_name("v-btn").click()

    # click to not to register auth methods
    chrome.find_element_by_class_name("red--text").click()


def login_test(test_email):
    # TODO
    return 


def random_email(prefix_length=10, suffix="robot.com"):
    letters = string.ascii_lowercase
    prefix = ''.join(random.choice(letters) for i in range(prefix_length))

    return f'{prefix}@{suffix}'


test_email = random_email()

register_user_test(test_email)
