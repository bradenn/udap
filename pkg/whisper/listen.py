#  Copyright (c) 2022 Braden Nicholson

import io
import os
import requests
import speech_recognition as sr
import tempfile
import whisper as w
from pydub import AudioSegment

temp_dir = tempfile.mkdtemp()
save_path = os.path.join(temp_dir, "temp.wav")


# This code is derived from https://github.com/mallorbc/whisper_mic

def main():
    model = "base.en"
    audio_model = w.load_model(model)
    # load the speech recognizer and set the initial energy threshold and pause threshold
    r = sr.Recognizer()
    r.energy_threshold = 300
    r.pause_threshold = 0.5
    r.dynamic_energy_threshold = False

    with sr.Microphone(sample_rate=16000) as source:
        while True:
            # get and save audio to wav file
            audio = r.listen(source)
            data = io.BytesIO(audio.get_wav_data())
            audio_clip = AudioSegment.from_file(data)
            audio_clip.export(save_path, format="wav")

            result = audio_model.transcribe(save_path, language='english', fp16=False)

            predicted_text = result["text"]
            prepared = predicted_text.lower()
            prepared = prepared.replace("at least", "atlas")
            print("You said:" + prepared)
            req = requests.post(url="http://10.0.1.2:5055/recognized", json={'text': prepared})


main()
