import asyncio
import cv2 as cv2
import dlib
import face_recognition
import json
import math
import numpy as np
import os
import os.path
import pickle
import websockets
from face_recognition.face_recognition_cli import image_files_in_folder
from sklearn import neighbors

ALLOWED_EXTENSIONS = {'png', 'jpg', 'jpeg'}

print("Using Cuda: {} ({} device)".format(dlib.DLIB_USE_CUDA, dlib.cuda.get_num_devices()))


def train(train_dir, model_save_path="trained_knn_model.clf", n_neighbors=None, knn_algo='ball_tree', verbose=False):
    encodings = []
    names = []

    # Loop through each person in the training set
    for class_dir in os.listdir(train_dir):
        if not os.path.isdir(os.path.join(train_dir, class_dir)):
            continue

        # Loop through each training image for the current person
        for img_path in image_files_in_folder(os.path.join(train_dir, class_dir)):
            image = face_recognition.load_image_file(img_path)
            face_bounding_boxes = face_recognition.face_locations(image, model="cnn")

            if len(face_bounding_boxes) != 1:
                # If there are no people (or too many people) in a training image, skip the image.
                if verbose:
                    print("Image {} not suitable for training: {}".format(img_path, "Didn't find a face" if len(
                        face_bounding_boxes) < 1 else "Found more than one face"))
            else:
                # Add face encoding for current image to the training set
                encodings.append(face_recognition.face_encodings(image, known_face_locations=face_bounding_boxes)[0])
                names.append(class_dir)

    # Determine how many neighbors to use for weighting in the KNN classifier
    if n_neighbors is None:
        n_neighbors = int(round(math.sqrt(len(encodings))))
        if verbose:
            print("Chose n_neighbors automatically:", n_neighbors)

    # Create and train the KNN classifier
    knn_clf = neighbors.KNeighborsClassifier(n_neighbors=n_neighbors, algorithm=knn_algo, weights='distance')
    knn_clf.fit(encodings, names)

    # Save the trained KNN classifier
    if model_save_path is not None:
        with open(model_save_path, 'wb') as f:
            pickle.dump(knn_clf, f)

    return knn_clf


def predict(frame, knn_clf=None, model_path=None, distance_threshold=0.5):
    if knn_clf is None and model_path is None:
        raise Exception("Must supply knn classifier either thourgh knn_clf or model_path")

    # Load a trained KNN model (if one was passed in)
    if knn_clf is None:
        with open(model_path, 'rb') as f:
            knn_clf = pickle.load(f)

    locations = face_recognition.face_locations(frame, model="hog")

    # If no faces are found in the image, return an empty result.
    if len(locations) == 0:
        return []

    # Find encodings for faces in the test image
    faces_encodings = face_recognition.face_encodings(frame, known_face_locations=locations)
    landmarks = face_recognition.face_landmarks(frame, face_locations=locations, model="small")
    # Use the KNN model to find the best matches for the test face
    closest_distances = knn_clf.kneighbors(faces_encodings, n_neighbors=1)
    are_matches = [closest_distances[0][i][0] <= distance_threshold for i in range(len(locations))]
    distances = [closest_distances[0][i][0] for i in range(len(locations))]

    # Predict classes and remove classifications that aren't within the threshold
    return [(pred, loc, landmarks, distances) if rec else ("unknown", loc, landmarks, distances) for
            pred, loc, rec, landmarks, distances in
            zip(knn_clf.predict(faces_encodings), locations, are_matches, landmarks, distances)]


async def recognize(websocket, path):
    global classifier

    async for message in websocket:
        nparr = np.frombuffer(message, dtype=np.uint8)
        img_np = cv2.imdecode(nparr, flags=1)
        predictions = predict(img_np, knn_clf=classifier)

        predicts = []

        for name, (top, right, bottom, left), landmarks, distances in predictions:
            predicts.append({
                "name": name,
                "top": top,
                "right": right,
                "bottom": bottom,
                "left": left,
                "distance": distances,
                "landmarks": {
                    "rightEye": {
                        "xa": landmarks["right_eye"][0][0],
                        "ya": landmarks["right_eye"][0][1],
                        "xb": landmarks["right_eye"][1][0],
                        "yb": landmarks["right_eye"][1][1],
                    },
                    "leftEye": {
                        "xa": landmarks["left_eye"][0][0],
                        "ya": landmarks["left_eye"][0][1],
                        "xb": landmarks["left_eye"][1][0],
                        "yb": landmarks["left_eye"][1][1],
                    },
                    "nose": {
                        "x": landmarks["nose_tip"][0][0],
                        "y": landmarks["nose_tip"][0][1],
                    }

                }
            })
        try:
            await websocket.send(json.dumps(predicts))
        except websockets.ConnectionClosed as e:
            continue


async def start():
    global classifier

    print("Training KNN classifier...")
    classifier = train("./known", model_save_path="trained_knn_model.clf", n_neighbors=2)
    print("Training complete!")

    async with websockets.serve(recognize, "0.0.0.0", 8765):
        await asyncio.Future()


if __name__ == '__main__':
    asyncio.run(start())
