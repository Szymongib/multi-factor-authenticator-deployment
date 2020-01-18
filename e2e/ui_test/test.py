from selenium import webdriver
import random
import string


def fill_credentials_input(driver, email, password):
    # login input id = login
    email_input = driver.find_element_by_id("input-12")
    email_input.clear()
    email_input.send_keys(email)

    # password input id = password
    password_input = driver.find_element_by_id("password")
    password_input.clear()
    password_input.send_keys(password)


def register_user_test(test_email):
    chrome.get('http://localhost:8001/')

    # Wait for any redirects
    chrome.implicitly_wait(5) # seconds

    # Login with Multi-Factor
    chrome.find_element_by_id("login-button").click()

    # Select: Register new user
    chrome.find_elements_by_xpath("//*[text()='Register new user']")[0].click()

    # Fill credentials
    fill_credentials_input(chrome, test_email, "secret-password")

    # click register button
    chrome.find_element_by_class_name("v-btn").click()

    chrome.implicitly_wait(5) # seconds

    # click to not to register auth methods
    chrome.find_element_by_class_name("red--text").click()

    chrome.implicitly_wait(10) # seconds


def login_test(test_email):
    chrome.get('http://localhost:8001/')

    # Wait for any redirects
    chrome.implicitly_wait(5) # seconds

    # Login with Multi-Factor
    chrome.find_element_by_id("login-button").click()

    chrome.implicitly_wait(5) # seconds

    fill_credentials_input(chrome, test_email, "secret-password")

    # click register button
    chrome.find_element_by_class_name("v-btn").click()

    chrome.implicitly_wait(10) # seconds


def add_todo(content):
    todo_input = chrome.find_element_by_id("itemInput")
    todo_input.clear()
    todo_input.send_keys(content)

    chrome.find_element_by_id("addButton").click()


def random_email(prefix_length=10, suffix="robot.com"):
    letters = string.ascii_lowercase
    prefix = ''.join(random.choice(letters) for i in range(prefix_length))

    return f'{prefix}@{suffix}'


chrome = webdriver.Chrome()

test_email = random_email()

register_user_test(test_email)
login_test(test_email)
add_todo("My robotic TODO")

