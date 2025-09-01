# 🚀 Deployment Guide

This guide will walk you through deploying your React app to Vercel for free!

## 📋 Prerequisites

- [Git](https://git-scm.com/) installed on your computer
- [GitHub](https://github.com/) account
- [Node.js](https://nodejs.org/) (version 14 or higher)
- Your React app code

## 🎯 Step 1: Prepare Your App

### 1.1 Clean Up Development Code
Your app has been prepared for production with:
- ✅ Production-safe logging (console.logs only show in development)
- ✅ Environment-based configuration
- ✅ Production build scripts
- ✅ Vercel deployment configuration

### 1.2 Test Production Build Locally
```bash
# Install dependencies (if not already done)
npm install

# Test production build
npm run build:prod

# Preview production build
npm run deploy:preview
```

## 🚀 Step 2: Deploy to Vercel

### 2.1 Push Code to GitHub
```bash
# Initialize git if not already done
git init

# Add all files
git add .

# Commit changes
git commit -m "Prepare for deployment"

# Add your GitHub repository as remote
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO_NAME.git

# Push to GitHub
git push -u origin main
```

### 2.2 Deploy with Vercel

1. **Go to [vercel.com](https://vercel.com)**
2. **Sign up/Login with GitHub**
3. **Click "New Project"**
4. **Import your GitHub repository**
5. **Vercel will automatically detect it's a React app**
6. **Configure environment variables:**
   ```
   REACT_APP_API_BASE_URL=https://your-backend-domain.com/api/v1
   REACT_APP_WS_BASE_URL=wss://your-backend-domain.com/api/v1
   REACT_APP_ENVIRONMENT=production
   ```
7. **Click "Deploy"**

## 🔧 Step 3: Configure Backend

### 3.1 Update API Endpoints
Your app is configured to use environment variables for API endpoints:

- **Development**: `http://localhost:8080/api/v1`
- **Production**: `https://your-backend-domain.com/api/v1`

### 3.2 Deploy Your Go Backend
You'll need to deploy your Go backend separately. Options include:
- **Railway** (free tier available)
- **Render** (free tier available)
- **DigitalOcean App Platform** (free tier available)
- **Heroku** (limited free tier)

## 🔄 Step 4: Update Live Site

### 4.1 Make Changes
1. **Edit your code locally**
2. **Test changes**
3. **Commit and push to GitHub:**
   ```bash
   git add .
   git commit -m "Add new feature"
   git push origin main
   ```
4. **Vercel automatically redeploys** your site!

### 4.2 Environment Variables
To update environment variables after deployment:
1. **Go to your Vercel project dashboard**
2. **Click "Settings" → "Environment Variables"**
3. **Add/Update variables**
4. **Redeploy** (automatic or manual)

## 🌐 Step 5: Custom Domain (Optional)

### 5.1 Add Custom Domain
1. **In Vercel dashboard, go to "Settings" → "Domains"**
2. **Add your custom domain**
3. **Follow DNS configuration instructions**
4. **Wait for DNS propagation (up to 48 hours)**

## 📱 Step 6: Test Everything

### 6.1 Test Your Deployed App
- ✅ **Homepage loads correctly**
- ✅ **User registration/login works**
- ✅ **Room creation works**
- ✅ **Video calling works**
- ✅ **All features function properly**

### 6.2 Monitor Performance
- **Vercel Analytics** (built-in)
- **Core Web Vitals** monitoring
- **Error tracking** (consider adding Sentry)

## 🚨 Troubleshooting

### Common Issues:

1. **Build Fails**
   - Check console for errors
   - Ensure all dependencies are installed
   - Verify environment variables

2. **API Calls Fail**
   - Check backend deployment
   - Verify API_BASE_URL environment variable
   - Check CORS configuration on backend

3. **Video Calling Issues**
   - Ensure WebSocket URLs are correct
   - Check backend WebSocket implementation
   - Verify SSL certificates for production

## 🎉 Success!

Your app is now deployed and will automatically update whenever you push changes to GitHub!

## 📚 Next Steps

- **Set up monitoring** (Sentry, LogRocket)
- **Add analytics** (Google Analytics, Mixpanel)
- **Implement CI/CD** (GitHub Actions)
- **Add testing** (Jest, Cypress)

## 🆘 Need Help?

- **Vercel Documentation**: [vercel.com/docs](https://vercel.com/docs)
- **React Deployment**: [reactjs.org/docs/deployment.html](https://reactjs.org/docs/deployment.html)
- **GitHub Issues**: Check your repository issues

---

**Happy Deploying! 🚀** 