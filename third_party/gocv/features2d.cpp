#include "features2d.h"

AKAZE AKAZE_Create() {
    try {
        return new cv::Ptr<cv::AKAZE>(cv::AKAZE::create());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}
AKAZE AKAZE_CreateWithParams(int descriptor_type, int descriptor_size, int descriptor_channels,
                             float threshold, int nOctaves, int nOctaveLayers, int diffusivity) {
    try {
        cv::AKAZE::DescriptorType type = static_cast<cv::AKAZE::DescriptorType>(descriptor_type);

        return new cv::Ptr<cv::AKAZE>(cv::AKAZE::create(type, descriptor_size, descriptor_channels,threshold, nOctaves, nOctaveLayers, static_cast<cv::KAZE::DiffusivityType>(diffusivity)));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}


void AKAZE_Close(AKAZE a) {
    delete a;
}

struct KeyPoints AKAZE_Detect(AKAZE a, Mat src) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*a)->detect(*src, detected);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

struct KeyPoints AKAZE_Compute(AKAZE a, Mat src, struct KeyPoints kp, Mat desc) {
    try {
        std::vector<cv::KeyPoint> computed;
        for (size_t i = 0; i < kp.length; i++) {
            cv::KeyPoint k = cv::KeyPoint(kp.keypoints[i].x, kp.keypoints[i].y,
                kp.keypoints[i].size, kp.keypoints[i].angle, kp.keypoints[i].response,
                kp.keypoints[i].octave, kp.keypoints[i].classID);
            computed.push_back(k);
        }
    
        (*a)->compute(*src, computed, *desc);
    
        KeyPoint* kps = new KeyPoint[computed.size()];
    
        for (size_t i = 0; i < computed.size(); ++i) {
            KeyPoint k = {computed[i].pt.x, computed[i].pt.y, computed[i].size, computed[i].angle,
                          computed[i].response, computed[i].octave, computed[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)computed.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

struct KeyPoints AKAZE_DetectAndCompute(AKAZE a, Mat src, Mat mask, Mat desc) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*a)->detectAndCompute(*src, *mask, detected, *desc);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

AgastFeatureDetector AgastFeatureDetector_Create() {
    try {
        return new cv::Ptr<cv::AgastFeatureDetector>(cv::AgastFeatureDetector::create());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

AgastFeatureDetector AgastFeatureDetector_CreateWithParams(int threshold, bool nonmaxSuppression, int type) {
    try {
        cv::AgastFeatureDetector::DetectorType detectorType = static_cast<cv::AgastFeatureDetector::DetectorType>(type);
        return new cv::Ptr<cv::AgastFeatureDetector>(cv::AgastFeatureDetector::create(threshold, nonmaxSuppression, detectorType));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

void AgastFeatureDetector_Close(AgastFeatureDetector a) {
    delete a;
}

struct KeyPoints AgastFeatureDetector_Detect(AgastFeatureDetector a, Mat src) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*a)->detect(*src, detected);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

BRISK BRISK_Create() {
    try {
        return new cv::Ptr<cv::BRISK>(cv::BRISK::create());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

BRISK BRISK_CreateWithParams(int thresh, int octaves, float patternScale) {
    try {
        return new cv::Ptr<cv::BRISK>(cv::BRISK::create(thresh, octaves, patternScale));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

void BRISK_Close(BRISK b) {
    delete b;
}

struct KeyPoints BRISK_Detect(BRISK b, Mat src) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*b)->detect(*src, detected);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

struct KeyPoints BRISK_Compute(BRISK b, Mat src, struct KeyPoints kp, Mat desc) {
    try {
        std::vector<cv::KeyPoint> computed;
        for (size_t i = 0; i < kp.length; i++) {
            cv::KeyPoint k = cv::KeyPoint(kp.keypoints[i].x, kp.keypoints[i].y,
                kp.keypoints[i].size, kp.keypoints[i].angle, kp.keypoints[i].response,
                kp.keypoints[i].octave, kp.keypoints[i].classID);
            computed.push_back(k);
        }
    
        (*b)->compute(*src, computed, *desc);
    
        KeyPoint* kps = new KeyPoint[computed.size()];
    
        for (size_t i = 0; i < computed.size(); ++i) {
            KeyPoint k = {computed[i].pt.x, computed[i].pt.y, computed[i].size, computed[i].angle,
                          computed[i].response, computed[i].octave, computed[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)computed.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

struct KeyPoints BRISK_DetectAndCompute(BRISK b, Mat src, Mat mask, Mat desc) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*b)->detectAndCompute(*src, *mask, detected, *desc);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;    
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

GFTTDetector GFTTDetector_Create() {
    try {
        return new cv::Ptr<cv::GFTTDetector>(cv::GFTTDetector::create());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

GFTTDetector GFTTDetector_Create_WithParams(const GFTTDetectorParams* params) {
    try {
        // Create the GFTTDetector and return it wrapped in a smart pointer
        return new cv::Ptr<cv::GFTTDetector>(cv::GFTTDetector::create(params->maxCorners, params->qualityLevel, params->minDistance,
            params->blockSize, params->useHarrisDetector, params->k));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}



void GFTTDetector_Close(GFTTDetector a) {
    delete a;
}

struct KeyPoints GFTTDetector_Detect(GFTTDetector a, Mat src) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*a)->detect(*src, detected);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

KAZE KAZE_Create() {
    try {
        return new cv::Ptr<cv::KAZE>(cv::KAZE::create());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

KAZE KAZE_CreateWithParams(bool extended, bool upright, float threshold, int nOctaves, int nOctaveLayers, int diffusivity) {
    try {
        return new cv::Ptr<cv::KAZE>(cv::KAZE::create(extended, upright, threshold, nOctaves, nOctaveLayers, static_cast<cv::KAZE::DiffusivityType>(diffusivity)));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

void KAZE_Close(KAZE a) {
    delete a;
}

struct KeyPoints KAZE_Detect(KAZE a, Mat src) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*a)->detect(*src, detected);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

struct KeyPoints KAZE_Compute(KAZE a, Mat src, struct KeyPoints kp, Mat desc) {
    try {
        std::vector<cv::KeyPoint> computed;
        for (size_t i = 0; i < kp.length; i++) {
            cv::KeyPoint k = cv::KeyPoint(kp.keypoints[i].x, kp.keypoints[i].y,
                kp.keypoints[i].size, kp.keypoints[i].angle, kp.keypoints[i].response,
                kp.keypoints[i].octave, kp.keypoints[i].classID);
            computed.push_back(k);
        }
    
        (*a)->compute(*src, computed, *desc);
    
        KeyPoint* kps = new KeyPoint[computed.size()];
    
        for (size_t i = 0; i < computed.size(); ++i) {
            KeyPoint k = {computed[i].pt.x, computed[i].pt.y, computed[i].size, computed[i].angle,
                          computed[i].response, computed[i].octave, computed[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)computed.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

struct KeyPoints KAZE_DetectAndCompute(KAZE a, Mat src, Mat mask, Mat desc) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*a)->detectAndCompute(*src, *mask, detected, *desc);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

MSER MSER_Create() {
    try {
        return new cv::Ptr<cv::MSER>(cv::MSER::create());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

MSER MSER_CreateWithParams(int delta, int min_area, int max_area, double max_variation, double min_diversity,
                 int max_evolution, double area_threshold, double min_margin, int edge_blur_size) {
    try {
        return new cv::Ptr<cv::MSER>(cv::MSER::create(delta, min_area, max_area, max_variation, min_diversity,
            max_evolution, area_threshold, min_margin, edge_blur_size));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }                
}

void MSER_Close(MSER a) {
    delete a;
}

struct KeyPoints MSER_Detect(MSER a, Mat src) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*a)->detect(*src, detected);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

FastFeatureDetector FastFeatureDetector_Create() {
    try {
        return new cv::Ptr<cv::FastFeatureDetector>(cv::FastFeatureDetector::create());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

void FastFeatureDetector_Close(FastFeatureDetector f) {
    delete f;
}

FastFeatureDetector FastFeatureDetector_CreateWithParams(int threshold, bool nonmaxSuppression, int type) {
    try {
        return new cv::Ptr<cv::FastFeatureDetector>(cv::FastFeatureDetector::create(threshold,nonmaxSuppression,static_cast<cv::FastFeatureDetector::DetectorType>(type)));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

struct KeyPoints FastFeatureDetector_Detect(FastFeatureDetector f, Mat src) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*f)->detect(*src, detected);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

ORB ORB_Create() {
    try {
        return new cv::Ptr<cv::ORB>(cv::ORB::create());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

ORB ORB_CreateWithParams(int nfeatures, float scaleFactor, int nlevels, int edgeThreshold, int firstLevel, int WTA_K, int scoreType, int patchSize, int fastThreshold) {
    try {
        return new cv::Ptr<cv::ORB>(cv::ORB::create(nfeatures, scaleFactor, nlevels, edgeThreshold, firstLevel, WTA_K, static_cast<cv::ORB::ScoreType>(scoreType), patchSize, fastThreshold));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

void ORB_Close(ORB o) {
    delete o;
}

struct KeyPoints ORB_Detect(ORB o, Mat src) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*o)->detect(*src, detected);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

struct KeyPoints ORB_Compute(ORB o, Mat src, struct KeyPoints kp, Mat desc) {
    try {
        std::vector<cv::KeyPoint> computed;
        for (size_t i = 0; i < kp.length; i++) {
            cv::KeyPoint k = cv::KeyPoint(kp.keypoints[i].x, kp.keypoints[i].y,
                kp.keypoints[i].size, kp.keypoints[i].angle, kp.keypoints[i].response,
                kp.keypoints[i].octave, kp.keypoints[i].classID);
            computed.push_back(k);
        }
    
        (*o)->compute(*src, computed, *desc);
    
        KeyPoint* kps = new KeyPoint[computed.size()];
    
        for (size_t i = 0; i < computed.size(); ++i) {
            KeyPoint k = {computed[i].pt.x, computed[i].pt.y, computed[i].size, computed[i].angle,
                          computed[i].response, computed[i].octave, computed[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)computed.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

struct KeyPoints ORB_DetectAndCompute(ORB o, Mat src, Mat mask, Mat desc) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*o)->detectAndCompute(*src, *mask, detected, *desc);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

cv::SimpleBlobDetector::Params ConvertCParamsToCPPParams(SimpleBlobDetectorParams params) {
    cv::SimpleBlobDetector::Params converted;

    converted.blobColor = params.blobColor;
    converted.filterByArea = params.filterByArea;
    converted.filterByCircularity = params.filterByCircularity;
    converted.filterByColor = params.filterByColor;
    converted.filterByConvexity = params.filterByConvexity;
    converted.filterByInertia = params.filterByInertia;
    converted.maxArea = params.maxArea;
    converted.maxCircularity = params.maxCircularity;
    converted.maxConvexity = params.maxConvexity;
    converted.maxInertiaRatio = params.maxInertiaRatio;
    converted.maxThreshold = params.maxThreshold;
    converted.minArea = params.minArea;
    converted.minCircularity = params.minCircularity;
    converted.minConvexity = params.minConvexity;
    converted.minDistBetweenBlobs = params.minDistBetweenBlobs;
    converted.minInertiaRatio = params.minInertiaRatio;
    converted.minRepeatability = params.minRepeatability;
    converted.minThreshold = params.minThreshold;
    converted.thresholdStep = params.thresholdStep;

    return converted;
}

SimpleBlobDetectorParams ConvertCPPParamsToCParams(cv::SimpleBlobDetector::Params params) {
    SimpleBlobDetectorParams converted;

    converted.blobColor = params.blobColor;
    converted.filterByArea = params.filterByArea;
    converted.filterByCircularity = params.filterByCircularity;
    converted.filterByColor = params.filterByColor;
    converted.filterByConvexity = params.filterByConvexity;
    converted.filterByInertia = params.filterByInertia;
    converted.maxArea = params.maxArea;
    converted.maxCircularity = params.maxCircularity;
    converted.maxConvexity = params.maxConvexity;
    converted.maxInertiaRatio = params.maxInertiaRatio;
    converted.maxThreshold = params.maxThreshold;
    converted.minArea = params.minArea;
    converted.minCircularity = params.minCircularity;
    converted.minConvexity = params.minConvexity;
    converted.minDistBetweenBlobs = params.minDistBetweenBlobs;
    converted.minInertiaRatio = params.minInertiaRatio;
    converted.minRepeatability = params.minRepeatability;
    converted.minThreshold = params.minThreshold;
    converted.thresholdStep = params.thresholdStep;

    return converted;
}

SimpleBlobDetector SimpleBlobDetector_Create_WithParams(SimpleBlobDetectorParams params){
    try {
        cv::SimpleBlobDetector::Params actualParams;
        return new cv::Ptr<cv::SimpleBlobDetector>(cv::SimpleBlobDetector::create(ConvertCParamsToCPPParams(params)));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

SimpleBlobDetector SimpleBlobDetector_Create() {
    try {
        return new cv::Ptr<cv::SimpleBlobDetector>(cv::SimpleBlobDetector::create());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

SimpleBlobDetectorParams SimpleBlobDetectorParams_Create() {
    try {
        return ConvertCPPParamsToCParams(cv::SimpleBlobDetector::Params());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return SimpleBlobDetectorParams();
    }
}

void SimpleBlobDetector_Close(SimpleBlobDetector b) {
    delete b;
}

struct KeyPoints SimpleBlobDetector_Detect(SimpleBlobDetector b, Mat src) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*b)->detect(*src, detected);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

BFMatcher BFMatcher_Create() {
    try {
        return new cv::Ptr<cv::BFMatcher>(cv::BFMatcher::create());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

BFMatcher BFMatcher_CreateWithParams(int normType, bool crossCheck) {
    try {
        return new cv::Ptr<cv::BFMatcher>(cv::BFMatcher::create(normType, crossCheck));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

void BFMatcher_Close(BFMatcher b) {
    delete b;
}

struct DMatches BFMatcher_Match(BFMatcher b, Mat query, Mat train) {
    try {
        std::vector<cv::DMatch> matches;
        (*b)->match(*query, *train, matches);
    
        DMatch *dmatches = new DMatch[matches.size()];
        for (size_t i = 0; i < matches.size(); ++i) {
            DMatch dmatch = {matches[i].queryIdx, matches[i].trainIdx, matches[i].imgIdx, matches[i].distance};
            dmatches[i] = dmatch;
        }
        DMatches ret = {dmatches, (int) matches.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        DMatch *dmatches = new DMatch[0];
        DMatches ret = {dmatches, 0};
        return ret;
    }
}

struct MultiDMatches BFMatcher_KnnMatch(BFMatcher b, Mat query, Mat train, int k) {
    try {
        std::vector< std::vector<cv::DMatch> > matches;
        (*b)->knnMatch(*query, *train, matches, k);
    
        DMatches *dms = new DMatches[matches.size()];
        for (size_t i = 0; i < matches.size(); ++i) {
            DMatch *dmatches = new DMatch[matches[i].size()];
            for (size_t j = 0; j < matches[i].size(); ++j) {
                DMatch dmatch = {matches[i][j].queryIdx, matches[i][j].trainIdx, matches[i][j].imgIdx,
                                 matches[i][j].distance};
                dmatches[j] = dmatch;
            }
            dms[i] = {dmatches, (int) matches[i].size()};
        }
        MultiDMatches ret = {dms, (int) matches.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        DMatch *dmatches = new DMatch[0];
        DMatches *dms = new DMatches[0];
        MultiDMatches ret = {dms, 0};
        return ret;
    }
}

struct MultiDMatches BFMatcher_KnnMatchWithParams(BFMatcher b, Mat query, Mat train, int k, Mat mask, bool compactResult) {
    try {
        std::vector< std::vector<cv::DMatch> > matches;
        (*b)->knnMatch(*query, *train, matches, k, *mask, compactResult);
    
        DMatches *dms = new DMatches[matches.size()];
        for (size_t i = 0; i < matches.size(); ++i) {
            DMatch *dmatches = new DMatch[matches[i].size()];
            for (size_t j = 0; j < matches[i].size(); ++j) {
                DMatch dmatch = {matches[i][j].queryIdx, matches[i][j].trainIdx, matches[i][j].imgIdx,
                                 matches[i][j].distance};
                dmatches[j] = dmatch;
            }
            dms[i] = {dmatches, (int) matches[i].size()};
        }
        MultiDMatches ret = {dms, (int) matches.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        DMatch *dmatches = new DMatch[0];
        DMatches *dms = new DMatches[0];
        MultiDMatches ret = {dms, 0};
        return ret;
    }
}

FlannBasedMatcher FlannBasedMatcher_Create() {
    try {
        return new cv::Ptr<cv::FlannBasedMatcher>(cv::FlannBasedMatcher::create());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

void FlannBasedMatcher_Close(FlannBasedMatcher f) {
    delete f;
}

struct MultiDMatches FlannBasedMatcher_KnnMatch(FlannBasedMatcher f, Mat query, Mat train, int k) {
    try {
        std::vector< std::vector<cv::DMatch> > matches;
        (*f)->knnMatch(*query, *train, matches, k);
    
        DMatches *dms = new DMatches[matches.size()];
        for (size_t i = 0; i < matches.size(); ++i) {
            DMatch *dmatches = new DMatch[matches[i].size()];
            for (size_t j = 0; j < matches[i].size(); ++j) {
                DMatch dmatch = {matches[i][j].queryIdx, matches[i][j].trainIdx, matches[i][j].imgIdx,
                                 matches[i][j].distance};
                dmatches[j] = dmatch;
            }
            dms[i] = {dmatches, (int) matches[i].size()};
        }
        MultiDMatches ret = {dms, (int) matches.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        DMatch *dmatches = new DMatch[0];
        DMatches *dms = new DMatches[0];
        MultiDMatches ret = {dms, 0};
        return ret;
    }
}

struct MultiDMatches FlannBasedMatcher_KnnMatchWithParams(FlannBasedMatcher f, Mat query, Mat train, int k, Mat mask, bool compactResult) {
    try {
        std::vector< std::vector<cv::DMatch> > matches;
        (*f)->knnMatch(*query, *train, matches, k, *mask, compactResult);
    
        DMatches *dms = new DMatches[matches.size()];
        for (size_t i = 0; i < matches.size(); ++i) {
            DMatch *dmatches = new DMatch[matches[i].size()];
            for (size_t j = 0; j < matches[i].size(); ++j) {
                DMatch dmatch = {matches[i][j].queryIdx, matches[i][j].trainIdx, matches[i][j].imgIdx,
                                 matches[i][j].distance};
                dmatches[j] = dmatch;
            }
            dms[i] = {dmatches, (int) matches[i].size()};
        }
        MultiDMatches ret = {dms, (int) matches.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        DMatch *dmatches = new DMatch[0];
        DMatches *dms = new DMatches[0];
        MultiDMatches ret = {dms, 0};
        return ret;
    }
}

void DrawKeyPoints(Mat src, struct KeyPoints kp, Mat dst, Scalar s, int flags) {
    try {
        std::vector<cv::KeyPoint> keypts;
        cv::KeyPoint keypt;
    
        for (int i = 0; i < kp.length; ++i) {
                keypt = cv::KeyPoint(kp.keypoints[i].x, kp.keypoints[i].y,
                                kp.keypoints[i].size, kp.keypoints[i].angle, kp.keypoints[i].response,
                                kp.keypoints[i].octave, kp.keypoints[i].classID);
                keypts.push_back(keypt);
        }
    
        cv::Scalar color = cv::Scalar(s.val1, s.val2, s.val3, s.val4);    
        cv::drawKeypoints(*src, keypts, *dst, color, static_cast<cv::DrawMatchesFlags>(flags));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
    }
}

SIFT SIFT_Create() {
    try {
        return new cv::Ptr<cv::SIFT>(cv::SIFT::create());
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}

SIFT SIFT_CreateWithParams(int nfeatures, int nOctaveLayers, double contrastThreshold, double edgeThreshold, double sigma) {
    try {
        return new cv::Ptr<cv::SIFT>(cv::SIFT::create(nfeatures, nOctaveLayers, contrastThreshold, edgeThreshold, sigma));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        return NULL;
    }
}


void SIFT_Close(SIFT d) {
    delete d;
}

struct KeyPoints SIFT_Detect(SIFT d, Mat src) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*d)->detect(*src, detected);
    
        KeyPoint* kps = new KeyPoint[detected.size()];
    
        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

struct KeyPoints SIFT_Compute(SIFT d, Mat src, struct KeyPoints kp, Mat desc) {
    try {
        std::vector<cv::KeyPoint> computed;
        for (size_t i = 0; i < kp.length; i++) {
            cv::KeyPoint k = cv::KeyPoint(kp.keypoints[i].x, kp.keypoints[i].y,
                kp.keypoints[i].size, kp.keypoints[i].angle, kp.keypoints[i].response,
                kp.keypoints[i].octave, kp.keypoints[i].classID);
            computed.push_back(k);
        }
    
        (*d)->compute(*src, computed, *desc);
    
        KeyPoint* kps = new KeyPoint[computed.size()];
    
        for (size_t i = 0; i < computed.size(); ++i) {
            KeyPoint k = {computed[i].pt.x, computed[i].pt.y, computed[i].size, computed[i].angle,
                          computed[i].response, computed[i].octave, computed[i].class_id
                         };
            kps[i] = k;
        }
    
        KeyPoints ret = {kps, (int)computed.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

struct KeyPoints SIFT_DetectAndCompute(SIFT d, Mat src, Mat mask, Mat desc) {
    try {
        std::vector<cv::KeyPoint> detected;
        (*d)->detectAndCompute(*src, *mask, detected, *desc);

        KeyPoint* kps = new KeyPoint[detected.size()];

        for (size_t i = 0; i < detected.size(); ++i) {
            KeyPoint k = {detected[i].pt.x, detected[i].pt.y, detected[i].size, detected[i].angle,
                          detected[i].response, detected[i].octave, detected[i].class_id
                         };
            kps[i] = k;
        }

        KeyPoints ret = {kps, (int)detected.size()};
        return ret;
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
        KeyPoint* kps = new KeyPoint[0];
        KeyPoints ret = {kps, 0};
        return ret;
    }
}

void DrawMatches(Mat img1, struct KeyPoints kp1, Mat img2, struct KeyPoints kp2, struct DMatches matches1to2, Mat outImg, const Scalar matchesColor, const Scalar pointColor, struct ByteArray matchesMask, int flags) {
    try {
        std::vector<cv::KeyPoint> kp1vec, kp2vec;
        cv::KeyPoint keypt;

        for (int i = 0; i < kp1.length; ++i) {
            keypt = cv::KeyPoint(kp1.keypoints[i].x, kp1.keypoints[i].y,
                                kp1.keypoints[i].size, kp1.keypoints[i].angle, kp1.keypoints[i].response,
                                kp1.keypoints[i].octave, kp1.keypoints[i].classID);
            kp1vec.push_back(keypt);
        }

        for (int i = 0; i < kp2.length; ++i) {
            keypt = cv::KeyPoint(kp2.keypoints[i].x, kp2.keypoints[i].y,
                                kp2.keypoints[i].size, kp2.keypoints[i].angle, kp2.keypoints[i].response,
                                kp2.keypoints[i].octave, kp2.keypoints[i].classID);
            kp2vec.push_back(keypt);
        }

        cv::Scalar cvmatchescolor = cv::Scalar(matchesColor.val1, matchesColor.val2, matchesColor.val3, matchesColor.val4);
        cv::Scalar cvpointcolor = cv::Scalar(pointColor.val1, pointColor.val2, pointColor.val3, pointColor.val4);

        std::vector<cv::DMatch> dmatchvec;
        cv::DMatch dm;

        for (int i = 0; i < matches1to2.length; i++) {
            dm = cv::DMatch(matches1to2.dmatches[i].queryIdx, matches1to2.dmatches[i].trainIdx,
                            matches1to2.dmatches[i].imgIdx, matches1to2.dmatches[i].distance);
            dmatchvec.push_back(dm);
        }

        std::vector<char> maskvec;

        for (int i = 0; i < matchesMask.length; i++) {
            maskvec.push_back(matchesMask.data[i]);
        }

        cv::drawMatches(*img1, kp1vec, *img2, kp2vec, dmatchvec, *outImg, cvmatchescolor, cvpointcolor, maskvec, static_cast<cv::DrawMatchesFlags>(flags));
    } catch(const cv::Exception& e){
        setExceptionInfo(e.code, e.what());
    }
}
