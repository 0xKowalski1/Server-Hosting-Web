# Structure

## Layout
### Navbar
#### Logo+Company Name -  ✓
#### Home - ✓
#### Games - ✓ 
#### Support - ✓
#### Profile when authed - 
#### Sign in/Getting Started when not authed - ✓
#### My Gameservers/Log Out button when authed -  ✓

### Footer


## No auth Required

### Home Page -
#### What | Hero Header, Call to action-
#### Why  | Introduction to company, why we are better -
#### How  | Features - 
#### Where | Locations -

### Games Page -
#### Header - ✓ 
#### Search/Filters - ✓ / x 
#### List of game cards - ✓ 
#### Dont see what you are looking for contact section - ✓
#### Hook up links - 

### Support -
#### Contact Form - 
#### FAQ -

### Get Started/Pricing/Checkout - 
## Guided

## Advanced

### Log In Page -
#### Continue with google -  ✓ 
#### Continue with discord - ✓
#### Styling - 

## Auth Required

### Payment form - 

### Profile -

### Gameservers Management -
#### Gameserver cards - 
##### Gameserver static info - 
##### Gameserver dynamic info -
##### Start/Stop buttons -
##### Control panel link - 
##### Deploy/archive buttons - 
##### Edit button - 

#### New Gameserver Form - 
##### Allow users to create gameserver - ✓
##### Game dropdown - ✓
##### Memory/Storage sliders - ✓
##### Min mem/storage for game dropdown - 

#### Edit Gameserver Form - 
##### Edit Name -
##### Edit mem/storage

#### Delete Gameserver Button - 

### Gameserver control panel - 
#### Dashboard -
#### Console/Log Viewer -
#### Config - 
#### Automated Tasks - 
#### File Browser -
#### Text Editor - 
#### Backups Manager - 
#### Support -

# TODO

## Convert Mockup to start of web app - ✓ 
### Decide on tech stack - Go, Echo, Templ, HTMX, Tailwind
### Set up layout - ✓
### Setup nav - ✓
### Move images into static - ✓
### Add Games list - ✓
### Add Auth Pages -  ✓
### Refactor - ✓
### Setup tailwind properly -  ✓
### Setup htmx properly - ✓
### Make nav correctly track path - ✓

## Auth -
### Google OAuth -  ✓
### Discord OAuth - ✓
### Log out button -  ✓
### Style Auth Page -
### Ensure auth is enforced -  ✓ 
### Ensure email is unique -
### Ensure correct errors are displayed -
### Put cookie store secret in config - 
### Logout shouldent always full reload / - 

## Toasts -
### Allow for basic errors with toasts - 

## 404 Page -

## Games Page
### Get Games from DB -  ✓
### Allow Search for games - ✓

# Unimportant

### Handle errors in games page - 
### Put fonts in build -
### Navbar is not mobile responsive -
### First route currently results in two requests, one for layout and one for content, should only be one request for SEO - ✓

# Bugs

## Log out button not showing as red -  ✓
## Using browser back on 404 does not return the full layout - ✓
## Build proccess does not correctly set env to production -
## Search bar 'x' is blue and hard to see -
## Everything shifts when you all the search results disappear (games, due to scrollbar disappearing) -
## Tailwind was purging classes (like bg-red-700) that I was using, switched to CDN for now -  ✓

# Before Relase
## Cors -
## Oauth Tokens -
## Cookies - 
